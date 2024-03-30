package aws

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
)

// newCertificateBundle creates certificateBundle based on given pem encoded full chain.
func newCertificateBundle(fullChain []byte) (*certificateBundle, error) {
	cert, chainRaw := pem.Decode(fullChain)
	if cert == nil {
		return nil, fmt.Errorf("no pem block found in the fullChain")
	}
	certificate, err := x509.ParseCertificate(cert.Bytes)
	if err != nil {
		return nil, err
	}

	return &certificateBundle{
		certBlock: cert,
		Cert:      certificate,
		ChainRaw:  chainRaw,
	}, nil
}

// certificateBundle represents a bundle consist of a parsed public cert and rest part of the chain.
type certificateBundle struct {
	certBlock *pem.Block
	Cert      *x509.Certificate // parsed client cert
	ChainRaw  []byte            // pem encoded cert chain, exclude the client cert itself
}

// ClientCertRaw returns the pem encoded client public cert, not including the rest part of the chain.
func (b *certificateBundle) ClientCertRaw() []byte {
	return pem.EncodeToMemory(b.certBlock)
}

// VerifySerialNumber checks whether the cert has the same serial number as given hex string.
//
// The input string should look like: '9f:7a:a7:f3:f6:2a' or omit the colons.
func (b *certificateBundle) VerifySerialNumber(serialNumber string) bool {
	decode, err := hex.DecodeString(strings.ReplaceAll(serialNumber, ":", ""))
	if err != nil {
		return false
	}
	sn := new(big.Int)
	sn.SetBytes(decode)
	return sn.Cmp(b.Cert.SerialNumber) == 0
}

// ContainsAllDomains checks whether the cert contains all domains in the given list.
// The check is accomplished by a char-by-char comparison.
func (b *certificateBundle) ContainsAllDomains(domains []string) bool {
	if len(domains) == 0 {
		return true
	}

	// shortcut for 1 element check
	if len(domains) == 1 {
		for _, name := range b.Cert.DNSNames {
			if name == domains[0] {
				return true
			}
		}
		return false
	}

	set := make(map[string]struct{}, len(b.Cert.DNSNames))
	for _, name := range b.Cert.DNSNames {
		set[name] = struct{}{}
	}
	for _, domain := range domains {
		if _, ex := set[domain]; !ex {
			return false
		}
	}
	return true
}

func (b *certificateBundle) VerifyHostname(hostname string) bool {
	return b.Cert.VerifyHostname(hostname) == nil
}

// VerifyHostnames verifies whether ALL hostnames are valid for the cert.
func (b *certificateBundle) VerifyHostnames(hostnames []string) bool {
	for _, hostname := range hostnames {
		if b.Cert.VerifyHostname(hostname) != nil {
			return false
		}
	}
	return true
}

// GetSerialNumberHexString returns the serial number of the cert in hex string without colons or
// leading zeros.
// e.g. 9f7aa7f3f62a
func (b *certificateBundle) GetSerialNumberHexString() string {
	return fmt.Sprintf("%x", b.Cert.SerialNumber)
}
