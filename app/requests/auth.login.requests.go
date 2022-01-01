package requests

import "github.com/gookit/validate"

type LoginForm struct {
	Email    string `json:"email" xml:"email" form:"email" validate:"required|email"`
	Password string `json:"password" xml:"password" form:"password" validate:"required"`
}

// Messages you can custom validator error messages.
func (f LoginForm) Messages() map[string]string {
	return validate.MS{
		"required":    "{field} is required.",
		"Email.email": "Format {field} in invalid.",
	}
}

// Translates you can custom field translates.
func (f LoginForm) Translates() map[string]string {
	return validate.MS{
		"Email":    "email",
		"Password": "password",
	}
}
