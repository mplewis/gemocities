package user

import (
	"crypto/sha512"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/martinlindhe/base36"
	"github.com/sethvargo/go-password/password"
)

// generatePassword generates a random 32-character alphanumeric password.
func generatePassword() (string, error) {
	pass, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		return "", fmt.Errorf("error generating password: %w", err)
	}
	return pass, nil
}

func HashCertificate(cert *x509.Certificate) CertificateHash {
	hash := sha512.Sum512(cert.Raw)
	return CertificateHash(strings.ToLower(base36.EncodeBytes(hash[:]))[:32])
}
