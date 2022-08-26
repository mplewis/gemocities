package user

import (
	"errors"
	"time"

	"github.com/mplewis/ez3"
)

type Manager struct {
	Store ez3.EZ3
}

type NewUserArgs struct {
	CertificateHash
	Email string
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

func (m *Manager) Create(args NewUserArgs) error {
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
	return m.Set(User{
		Created:         time.Now(),
		EmailVerified:   false,
		CertificateHash: args.CertificateHash,
		Email:           args.Email,
		WebDAVPassword:  password,
	})
}

func (m *Manager) ChangePassword(user User) error {
	password, err := generatePassword()
	if err != nil {
		return err
	}
	user.WebDAVPassword = password
	return m.Set(user)
}

func (m *Manager) Verify(user User) error {
	user.EmailVerified = true
	return m.Set(user)
}
