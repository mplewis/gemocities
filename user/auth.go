package user

import (
	"errors"
	"fmt"
	"net/http"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/rs/zerolog/log"
)

type UserInfo struct {
	HasCertificate  bool
	CertificateHash string
	HasCommonName   bool
	CommonName      string
	HasUser         bool
	User            User
}

var ErrCertMismatch = errors.New("certificate mismatch")

func (m *Manager) AuthorizeWebDAVUser(r *http.Request) (bool, string) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false, ""
	}

	user, found, err := m.Get(username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return false, ""
	}
	if !found {
		return false, ""
	}

	if password != user.WebDAVPassword {
		return false, ""
	}
	return true, username
}

func (m *Manager) AuthorizeGeminiUser(r *gemini.Request) (UserInfo, error) {
	info := UserInfo{}

	tls := r.TLS()
	if len(tls.PeerCertificates) == 0 {
		return info, nil
	}
	cert := tls.PeerCertificates[0]
	info.HasCertificate = true
	info.CertificateHash = HashCertificate(cert)

	info.CommonName = cert.Subject.CommonName
	if info.CommonName == "" {
		return info, nil
	}
	info.HasCommonName = true

	user, found, err := m.Get(info.CommonName)
	if err != nil {
		return info, err
	}
	if !found {
		return info, nil
	}
	info.HasUser = true
	info.User = *user

	actual := HashCertificate(cert)
	if HashCertificate(cert) != user.CertificateHash {
		return info, fmt.Errorf("[%w] user: %s, expected %s, got: %s",
			ErrCertMismatch, user.CommonName, user.CertificateHash, actual)
	}
	return info, nil
}
