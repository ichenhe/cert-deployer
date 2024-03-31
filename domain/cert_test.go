package domain

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// full_chain.pem for testing.
//   - Name: *.chenhe.me
//   - Alternative name: [*.chenhe.me, chenhe.me]
//   - Serial number: 9f7aa7f3f62a992d9364d7f5f47b52b1
//   - Algorithm: SHA384withECDSA
//
//go:embed test_full_chain.pem
var testCert []byte

func Test_defaultCertificateBundle_GetRawCert(t *testing.T) {
	b, err := NewCertificateBundle(testCert)
	assert.NoError(t, err)

	// extract the first certificate from the full chain
	certStr := ""
	for _, s := range strings.Split(string(testCert), "\n") {
		certStr += s + "\n"
		if s == "-----END CERTIFICATE-----" {
			break
		}
	}

	assert.Equal(t, certStr, string(b.GetRawCert()))
}

func Test_defaultCertificateBundle_GetRaw(t *testing.T) {
	b, err := NewCertificateBundle(testCert)
	assert.NoError(t, err)

	want := string(testCert)
	actual := string(b.GetRaw())
	assert.Equal(t, want, actual)
}

func Test_defaultCertificateBundle_VerifySerialNumber(t *testing.T) {
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
			b, err := NewCertificateBundle(testCert)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, b.VerifySerialNumber(tt.serialNumber), "VerifySerialNumber(%v)", tt.serialNumber)
		})
	}
}

func Test_defaultCertificateBundle_ContainsAllDomains(t *testing.T) {
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
			b, err := NewCertificateBundle(testCert)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, b.ContainsAllDomains(tt.domains), "ContainsAllDomains(%v)", tt.domains)
		})
	}
}

func Test_defaultCertificateBundle_VerifyHostname(t *testing.T) {
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
			b, err := NewCertificateBundle(testCert)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, b.VerifyHostname(tt.hostname), "VerifyHostname(%v)", tt.hostname)
		})
	}
}

func Test_defaultCertificateBundle_VerifyHostnames(t *testing.T) {
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
			b, err := NewCertificateBundle(testCert)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, b.VerifyHostnames(tt.hostnames), "VerifyHostnames(%v)", tt.hostnames)
		})
	}
}

func Test_defaultCertificateBundle_GetSerialNumberHexString(t *testing.T) {
	want := "9f7aa7f3f62a992d9364d7f5f47b52b1"
	b, err := NewCertificateBundle(testCert)
	assert.NoError(t, err)
	assert.Equalf(t, want, b.GetSerialNumberHexString(), "GetSerialNumberHexString()")
}

func Test_defaultCertificateBundle_GetDomains(t *testing.T) {
	want := []string{"*.chenhe.me", "chenhe.me"}
	b, err := NewCertificateBundle(testCert)
	assert.NoError(t, err)
	assert.Equalf(t, want, b.GetDomains(), "GetDomains()")
}
