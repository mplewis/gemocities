package mail

import (
	"fmt"

	"github.com/mplewis/gemocities/user"
)

func (m Mailer) SendVerificationEmail(user user.User) error {
	return m.Send(Args{
		From:     fmt.Sprintf("Gemocities <welcome@%s>", m.AppDomain),
		To:       []string{user.Email},
		Subject:  "Confirm your Gemocities account",
		Template: "verify",
		Data:     user,
	})
}
