package user

import (
	"crypto/x509"
	"errors"
	"net/http"

	"git.sr.ht/~adnano/go-gemini"
)

type UserInfo struct {
	HasCertificate  bool
	CertificateHash CertificateHash
	HasUser         bool
	User            User
}

var ErrCertMismatch = errors.New("certificate mismatch")
var ErrNoUsername = errors.New("user has no username")

func (m *Manager) AuthorizeWebDAVUser(r *http.Request) (bool, User, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false, User{}, nil
	}

	user, found, err := m.Get(CertificateHash(username))
	if err != nil {
		return false, User{}, err
	}
	if !found {
		return false, User{}, nil
	}
	if !user.EmailVerified {
		return false, User{}, nil
	}
	if user.Name == "" {
		return false, User{}, ErrNoUsername
	}

	if password != user.WebDAVPassword {
		return false, User{}, nil
	}
	return true, user, nil
}

func (m *Manager) AuthorizeGeminiUser(r *gemini.Request) (UserInfo, error) {
	info := UserInfo{}

	cert := m.getCert(r)
	if cert == nil {
		return info, nil
	}
	info.HasCertificate = true
	info.CertificateHash = HashCertificate(cert)

	user, found, err := m.Get(info.CertificateHash)
	if err != nil {
		return info, err
	}
	if !found {
		return info, nil
	}

	info.HasUser = true
	info.User = user
	return info, nil
}

func (m *Manager) getCert(r *gemini.Request) *x509.Certificate {
	tls := r.TLS()
	if tls != nil && len(tls.PeerCertificates) != 0 {
		return tls.PeerCertificates[0]
	}

	if m.TestMode {
		// HACK: enables injection of certs during tests
		if r.Certificate != nil {
			return r.Certificate.Leaf
		}
	}

	return nil
}
