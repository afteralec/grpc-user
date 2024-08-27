package username

import "github.com/go-playground/validator/v10"

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func IsValid(username string) error {
	if err := validate.Var(username, "required,min=4,max=16,alphanum,lowercase"); err != nil {
		return err
	}
	return nil
}
