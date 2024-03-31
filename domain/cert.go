package domain

import (
	"bytes"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
)

// CertificateBundle represents a bundle consist of a parsed public cert and rest part of the chain.
//
// It was abstracted to an interface for testing only.
type CertificateBundle interface {
	// GetRawCert returns the pem encoded client public cert, not including the rest part of the chain.
	GetRawCert() []byte

	// GetRaw returns the pem encoded public cert, including the chain.
	GetRaw() []byte

	// GetDomains returns all DNS names in the cert.
	GetDomains() []string

	// VerifySerialNumber checks whether the cert has the same serial number as given hex string.
	//
	// The input string should look like: '9f:7a:a7:f3:f6:2a' or omit the colons.
	VerifySerialNumber(serialNumber string) bool

	// ContainsAllDomains checks whether the cert contains all domains in the given list.
	// The check is accomplished by a char-by-char comparison.
	ContainsAllDomains(domains []string) bool

	VerifyHostname(hostname string) bool

	// VerifyHostnames verifies whether ALL hostnames are valid for the cert.
	VerifyHostnames(hostnames []string) bool

	// GetSerialNumberHexString returns the serial number of the cert in hex string without colons or
	// leading zeros.
	// e.g. 9f7aa7f3f62a
	GetSerialNumberHexString() string
}

// NewCertificateBundle creates CertificateBundle based on given pem encoded full chain.
func NewCertificateBundle(fullChain []byte) (CertificateBundle, error) {
	cert, chainRaw := pem.Decode(fullChain)
	if cert == nil {
		return nil, fmt.Errorf("no pem block found in the fullChain")
	}
	certificate, err := x509.ParseCertificate(cert.Bytes)
	if err != nil {
		return nil, err
	}

	return &defaultCertificateBundle{
		certBlock: cert,
		cert:      certificate,
		ChainRaw:  chainRaw,
	}, nil
}

type defaultCertificateBundle struct {
	certBlock *pem.Block
	cert      *x509.Certificate // parsed client cert
	ChainRaw  []byte            // pem encoded cert chain, exclude the client cert itself
}

func (b *defaultCertificateBundle) GetDomains() []string {
	return b.cert.DNSNames
}

func (b *defaultCertificateBundle) GetRawCert() []byte {
	return pem.EncodeToMemory(b.certBlock)
}

func (b *defaultCertificateBundle) GetRaw() []byte {
	return bytes.Join([][]byte{b.GetRawCert(), b.ChainRaw}, []byte{})
}

func (b *defaultCertificateBundle) VerifySerialNumber(serialNumber string) bool {
	decode, err := hex.DecodeString(strings.ReplaceAll(serialNumber, ":", ""))
	if err != nil {
		return false
	}
	sn := new(big.Int)
	sn.SetBytes(decode)
	return sn.Cmp(b.cert.SerialNumber) == 0
}

// ContainsAllDomains checks whether the cert contains all domains in the given list.
// The check is accomplished by a char-by-char comparison.
func (b *defaultCertificateBundle) ContainsAllDomains(domains []string) bool {
	if len(domains) == 0 {
		return true
	}

	// shortcut for 1 element check
	if len(domains) == 1 {
		for _, name := range b.cert.DNSNames {
			if name == domains[0] {
				return true
			}
		}
		return false
	}

	set := make(map[string]struct{}, len(b.cert.DNSNames))
	for _, name := range b.cert.DNSNames {
		set[name] = struct{}{}
	}
	for _, domain := range domains {
		if _, ex := set[domain]; !ex {
			return false
		}
	}
	return true
}

func (b *defaultCertificateBundle) VerifyHostname(hostname string) bool {
	return b.cert.VerifyHostname(hostname) == nil
}

func (b *defaultCertificateBundle) VerifyHostnames(hostnames []string) bool {
	for _, hostname := range hostnames {
		if b.cert.VerifyHostname(hostname) != nil {
			return false
		}
	}
	return true
}

func (b *defaultCertificateBundle) GetSerialNumberHexString() string {
	return fmt.Sprintf("%x", b.cert.SerialNumber)
}
