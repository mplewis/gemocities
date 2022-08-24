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
	CommonName      string
	CertificateHash string
	Email           string
}

func (m *Manager) Get(CommonName string) (*User, bool, error) {
	user := &User{}
	err := m.Store.Get(CommonName, user)
	if errors.Is(err, ez3.KeyNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return user, true, nil
}

func (m *Manager) Set(user *User) error {
	return m.Store.Set(user.CertificateHash, user)
}

func (m *Manager) Create(args NewUserArgs) error {
	_, found, err := m.Get(args.CommonName)
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
	return m.Set(&User{
		Created:         time.Now(),
		EmailVerified:   false,
		CommonName:      args.CommonName,
		CertificateHash: args.CertificateHash,
		Email:           args.Email,
		WebDAVPassword:  password,
	})
}

func (m *Manager) ChangePassword(user *User) error {
	password, err := generatePassword()
	if err != nil {
		return err
	}
	user.WebDAVPassword = password
	return m.Set(user)
}
