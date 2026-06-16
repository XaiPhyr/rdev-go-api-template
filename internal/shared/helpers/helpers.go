package helpers

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func CleanSpecialChars(s string) string {
	re := regexp.MustCompile(`[!@#$%^&*()\_\+\=]`)
	clean := re.ReplaceAll([]byte(s), []byte(``))

	return string(clean)
}

func ValidateStruct(s any) error {
	v := validator.New()

	return v.Struct(s)
}

func ParseValidationErr(err error) map[string]string {
	var fields = make(map[string]string)

	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				fields[e.Field()] = ValidationErrMsg[e.StructNamespace()]
			}
		}
	}

	return fields
}
