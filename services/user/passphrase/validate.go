package passphrase

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	MinLength = 8
	MaxLength = 255
)

var validate *validator.Validate

var passphraseRegex = regexp.MustCompile(`^[a-zA-Z0-9_\!\@\#\$\%\^\&\* ]+$`)

func ValidatePassphrase(fl validator.FieldLevel) bool {
	return passphraseRegex.MatchString(fl.Field().String())
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("passphrase", ValidatePassphrase)
}

func IsValid(pw string) bool {
	if err := validate.Var(pw, "required,min=8,max=255,passphrase"); err != nil {
		return false
	}
	return true
}
