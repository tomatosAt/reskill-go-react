package util

import (
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

var errorMessage = map[string]string{
	"required": "The {field} field is required.",
	"email":    "The {field} field must be a valid email address.",
	"max":      "The {field} field must be at most {param} characters.",
	"min":      "The {field} field must be at least {param} characters.",
	"digits":   "The {field} must be {param} digits.",
	"len":      "The {field} must be {param} digits.",
	"number":   "The {field} must be {param} digits.",
}

type PHPValidatorResponse struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

func regexValidation(fl validator.FieldLevel) bool {
	regex := fl.Param()
	value := fl.Field().String()

	re := regexp.MustCompile(regex)
	return re.MatchString(value)
}

func getJSONFieldName(structType reflect.Type, fieldName string) string {
	if field, found := structType.Elem().FieldByName(fieldName); found {
		return field.Tag.Get("json")
	}
	return fieldName

}

func formatErrorMessage(tag string, field string, param string) string {
	tmpl, found := errorMessage[tag]
	if !found {
		tmpl = "Validation failed on the {field} field."
	}

	msg := strings.ReplaceAll(tmpl, "{field}", strings.ReplaceAll(field, "_", " "))
	msg = strings.ReplaceAll(msg, "{param}", param)

	return msg
}

func ValidatorStruct(payload interface{}) (string, error) {
	var validate = validator.New()
	validate.RegisterValidation("regexp", Regexp)
	if err := validate.Struct(payload); err != nil {
		logrus.Errorln("[ValidatorStruct] validate struct err -> ", err)
		errType := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			if _, ok := errType[err.Tag()]; !ok {
				errType[err.Tag()] += err.StructNamespace()
			} else {
				errType[err.Tag()] += ", " + err.StructNamespace()
			}
		}
		msg := ""
		i := 0
		for e, m := range errType {
			if e == "required" {
				msg += "[missing parameters]: " + m
			} else if e == "regexp" {
				msg += "[invalid format]: " + m
			} else {
				msg += e + ": " + m
			}
			if i < len(errType)-1 {
				msg += " | "
			}
			i++
		}
		return msg, err
	}
	return "", nil
}

func CheckInvalidId(id string) bool {
	emptyValues := []string{"", "#N/A", "-"}
	return slices.Contains(emptyValues, id)
}
