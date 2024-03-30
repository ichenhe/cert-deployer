package aws

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// fullchain.pem for testing.
//   - Name: *.chenhe.me
//   - Alternative name: [*.chenhe.me, chenhe.me]
//   - Serial number: 9f7aa7f3f62a992d9364d7f5f47b52b1
//   - Algorithm: SHA384withECDSA
//
//go:embed test_fullchain.pem
var testCert []byte

func Test_certificateBundle_ClientCertRaw(t *testing.T) {
	b, err := newCertificateBundle(testCert)
	assert.NoError(t, err)

	// extract the first certificate from the full chain
	certStr := ""
	for _, s := range strings.Split(string(testCert), "\n") {
		certStr += s + "\n"
		if s == "-----END CERTIFICATE-----" {
			break
		}
	}

	assert.Equal(t, certStr, string(b.ClientCertRaw()))
}

func Test_certificateBundle_VerifySerialNumber(t *testing.T) {
	tests := []struct {
		name         string
		serialNumber string
		want         bool
	}{
		{"matched", "9f7aa7f3f62a992d9364d7f5f47b52b1", true},
		{"matched with colons", "9f:7a:a7:f3:f6:2a:99:2d:93:64:d7:f5:f4:7b:52:b1", true},
		{"matched with leading zero", "009f7aa7f3f62a992d9364d7f5f47b52b1", true},
		{"not equal", "aa:bb:aa:bb:ff:aa:bb:aa:bb", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := newCertificateBundle(testCert)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, b.VerifySerialNumber(tt.serialNumber), "VerifySerialNumber(%v)", tt.serialNumber)
		})
	}
}

func Test_certificateBundle_ContainsAllDomains(t *testing.T) {
	tests := []struct {
		name    string
		domains []string
		want    bool
	}{
		{"empty", []string{}, true},
		{"all contained", []string{"*.chenhe.me", "chenhe.me"}, true},
		{"reversed order", []string{"chenhe.me", "*.chenhe.me"}, true},
		{"only logical matched", []string{"www.chenhe.me"}, false},
		{"not contained", []string{"chenhe.me", "chenhe.cc"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := newCertificateBundle(testCert)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, b.ContainsAllDomains(tt.domains), "ContainsAllDomains(%v)", tt.domains)
		})
	}
}

func Test_certificateBundle_VerifyHostname(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		want     bool
	}{
		{"matched", "chenhe.me", true},
		{"logical matched", "prefix.chenhe.me", true},
		{"not matched", "chenhe.cc", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := newCertificateBundle(testCert)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, b.VerifyHostname(tt.hostname), "VerifyHostname(%v)", tt.hostname)
		})
	}
}

func Test_certificateBundle_VerifyHostnames(t *testing.T) {
	tests := []struct {
		name      string
		hostnames []string
		want      bool
	}{
		{"all matched", []string{"chenhe.me", "*.chenhe.me"}, true},
		{"empty", []string{}, true},
		{"partial matched", []string{"chenhe.me", "chenhe.cc"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := newCertificateBundle(testCert)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, b.VerifyHostnames(tt.hostnames), "VerifyHostnames(%v)", tt.hostnames)
		})
	}
}
