package emailservice

import "errors"

type EmailDTO struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (d EmailDTO) Valid() error {
	if d.To == "" {
		return errors.New("email.To field cannot be empty")
	}

	if d.Subject == "" {
		return errors.New("email subject cannot be empty")
	}

	return nil
}
