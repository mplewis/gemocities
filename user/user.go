package user

import (
	"encoding/json"
	"time"
)

type User struct {
	Created         time.Time `json:"created"`
	CommonName      string    `json:"common_name"`
	CertificateHash string    `json:"certificate_hash"`
	Email           string    `json:"email"`
	EmailVerified   bool      `json:"email_verified"`
	WebDAVPassword  string    `json:"webdav_password"`
}

func (u *User) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) Deserialize(data []byte) error {
	return json.Unmarshal(data, u)
}
