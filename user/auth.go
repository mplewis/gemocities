package user

import (
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

func (m *Manager) AuthorizeWebDAVUser(r *http.Request) (bool, string, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false, "", nil
	}

	user, found, err := m.Get(CertificateHash(username))
	if err != nil {
		return false, "", err
	}
	if !found {
		return false, "", nil
	}

	if password != user.WebDAVPassword {
		return false, "", nil
	}
	return true, username, nil
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

	user, found, err := m.Get(info.CertificateHash)
	if err != nil {
		return info, err
	}
	if !found {
		return info, nil
	}

	info.HasUser = true
	info.User = *user
	return info, nil
}
