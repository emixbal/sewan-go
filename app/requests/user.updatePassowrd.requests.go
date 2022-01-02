package requests

import "github.com/gookit/validate"

type UserUpdatePasswordForm struct {
	Password string `json:"password" xml:"password" form:"password" validate:"required"`
}

// Messages you can custom validator error messages.
func (f UserUpdatePasswordForm) Messages() map[string]string {
	return validate.MS{
		"required": "{field} is required.",
	}
}

// Translates you can custom field translates.
func (f UserUpdatePasswordForm) Translates() map[string]string {
	return validate.MS{
		"Password": "password",
	}
}
