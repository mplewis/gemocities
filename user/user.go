package user

import (
	"encoding/json"
	"time"
)

type CertificateHash string

type User struct {
	Created           time.Time       `json:"created"`
	CertificateHash   CertificateHash `json:"certificate_hash"`
	Email             string          `json:"email"`
	EmailVerified     bool            `json:"email_verified"`
	VerificationToken string          `json:"verification_token"`
	Name              string          `json:"name"`
	WebDAVPassword    string          `json:"webdav_password"`
}

func (u *User) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) Deserialize(data []byte) error {
	return json.Unmarshal(data, u)
}
