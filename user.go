package gemocities

import (
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/mplewis/ez3"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-password/password"
)

type User struct {
	Created         time.Time `json:"created"`
	CommonName      string    `json:"common_name"`
	CertificateHash string    `json:"certificate_hash"`
	Email           string    `json:"email"`
	EmailVerified   bool      `json:"email_verified"`
	WebDAVPassword  string    `json:"webdav_password"`
}

type NewUserArgs struct {
	CommonName      string
	CertificateHash string
	Email           string
}

type GeminiError struct {
	gemini.Status
	Message string
}

func (u *User) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) Deserialize(data []byte) error {
	return json.Unmarshal(data, u)
}

type Users struct {
	Store ez3.EZ3
}

func (u Users) Get(CommonName string) (*User, bool, error) {
	user := &User{}
	err := u.Store.Get(CommonName, user)
	if errors.Is(err, ez3.KeyNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return user, true, nil
}

func (u Users) Set(user *User) error {
	return u.Store.Set(user.CertificateHash, user)
}

func (u Users) Create(args NewUserArgs) error {
	_, found, err := u.Get(args.CommonName)
	if err != nil {
		return err
	}
	if found {
		return errors.New("user already exists")
	}
	password, err := generatePassword()
	if err != nil {
		return err
	}
	return u.Set(&User{
		Created:         time.Now(),
		EmailVerified:   false,
		CommonName:      args.CommonName,
		CertificateHash: args.CertificateHash,
		Email:           args.Email,
		WebDAVPassword:  password,
	})
}

func (u Users) AuthorizeWebDAVUser(r *http.Request) (bool, string) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false, ""
	}

	user, found, err := u.Get(username)
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

func (u Users) AuthorizeGeminiUser(r *gemini.Request) (*User, bool, *GeminiError) {
	tls := r.TLS()
	if len(tls.PeerCertificates) == 0 {
		return nil, false, &GeminiError{gemini.StatusBadRequest, "client certificate required"}
	}

	cert := tls.PeerCertificates[0]
	cn := cert.Subject.CommonName
	if cn == "" {
		return nil, false, &GeminiError{gemini.StatusBadRequest, "client certificate lacks Common Name"}
	}

	user, found, err := u.Get(cn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return nil, false, &GeminiError{gemini.StatusTemporaryFailure, "internal server error"}
	}
	if !found {
		return nil, false, &GeminiError{gemini.StatusBadRequest, "user not found"}
	}

	if hashCertificate(cert) != user.CertificateHash {
		return nil, false, &GeminiError{gemini.StatusBadRequest, "user doesn't match client certificate"}
	}
	return user, true, nil
}

// generatePassword generates a random 32-character alphanumeric password.
func generatePassword() (string, error) {
	return password.Generate(32, 10, 0, false, false)
}

func hashCertificate(cert *x509.Certificate) string {
	hash := sha512.Sum512(cert.Raw)
	return fmt.Sprintf("%x", hash)
}
