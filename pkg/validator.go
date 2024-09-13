package pkg

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field    string      `json:"field"`
	Tag      string      `json:"message"`
	Value    interface{} `json:"-"`
	HasError bool        `json:"-"`
}

var validate = validator.New()

func ValidateStuct(data interface{}) (valid []ValidationError, mes string) {
	var validationError []ValidationError

	errs := validate.Struct(data)

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var val ValidationError

			val.Field = err.Field()
			val.Value = err.Value()
			val.Tag = err.Tag()
			val.HasError = true

			validationError = append(validationError, val)
		}
	}

	if errs := validationError; len(errs) > 0 && errs[0].HasError {
		fmt.Println(errs, "Array erros")
		errorMessage := make([]string, 0)
		for _, er := range errs {
			errorMessage = append(errorMessage, fmt.Sprintf("%s field has failed. Validation is: %s", er.Field, er.Tag))
		}

		return validationError, strings.Join(errorMessage, " and ")
	}

	return validationError, "done"
}
