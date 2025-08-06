package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#\$%\^&\*]{8,}$`)

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return passwordRegex.MatchString(password)
}
