package requests

import "github.com/gookit/validate"

type RegisterForm struct {
	Email    string `json:"email" xml:"email" form:"email" validate:"required|email"`
	Password string `json:"password" xml:"password" form:"password" validate:"required"`
	Name     string `json:"name" xml:"string" form:"name" validate:"required"`
}

// Messages you can custom validator error messages.
func (f RegisterForm) Messages() map[string]string {
	return validate.MS{
		"required":    "{field} is required.",
		"Email.email": "Format {field} in invalid.",
	}
}

// Translates you can custom field translates.
func (f RegisterForm) Translates() map[string]string {
	return validate.MS{
		"Email":    "email",
		"Name":     "name",
		"Password": "password",
	}
}
