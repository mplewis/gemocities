package user

import (
	"crypto/sha512"
	"crypto/x509"
	"strings"

	"github.com/martinlindhe/base36"
	"github.com/sethvargo/go-password/password"
)

// generatePassword generates a random 32-character alphanumeric password.
func generatePassword() (string, error) {
	return password.Generate(32, 10, 0, false, false)
}

func HashCertificate(cert *x509.Certificate) CertificateHash {
	hash := sha512.Sum512(cert.Raw)
	return CertificateHash(strings.ToLower(base36.EncodeBytes(hash[:]))[:32])
}
