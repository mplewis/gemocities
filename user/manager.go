package user

import (
	"errors"
	"time"

	"github.com/mplewis/ez3"
)

var ErrInvalidToken = errors.New("invalid token")

type Manager struct {
	Store ez3.EZ3
	Mailer
	TestMode bool // enables an alternate client cert parsing path
}

type NewArgs struct {
	CertificateHash
	Email    string
	Username string
}

type Mailer interface {
	SendVerificationEmail(user User) error
}

func (m *Manager) Get(ch CertificateHash) (User, bool, error) {
	user := &User{}
	err := m.Store.Get(string(ch), user)
	if errors.Is(err, ez3.KeyNotFound) {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, err
	}
	return *user, true, nil
}

func (m *Manager) Set(user User) error {
	return m.Store.Set(string(user.CertificateHash), &user)
}

func (m *Manager) Create(args NewArgs) error {
	_, found, err := m.Get(args.CertificateHash)
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
	token, err := generatePassword()
	if err != nil {
		return err
	}
	user := User{
		Created:           time.Now(),
		Name:              args.Username,
		EmailVerified:     false,
		CertificateHash:   args.CertificateHash,
		Email:             args.Email,
		WebDAVPassword:    password,
		VerificationToken: token,
	}
	err = m.Set(user)
	if err != nil {
		return err
	}
	// TODO
	// err = m.Mailer.SendVerificationEmail(user)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (m *Manager) ChangePassword(user User) error {
	password, err := generatePassword()
	if err != nil {
		return err
	}
	user.WebDAVPassword = password
	return m.Set(user)
}

func (m *Manager) Verify(user User, token string) error {
	if user.VerificationToken != token {
		return ErrInvalidToken
	}
	user.EmailVerified = true
	return m.Set(user)
}
