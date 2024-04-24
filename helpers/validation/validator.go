package validation

import (
	"addressBook/models/user"

	"github.com/go-playground/validator/v10"
)

func ValidateUserInput(input user.User) error {

	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		return err
	}
	return nil

}
