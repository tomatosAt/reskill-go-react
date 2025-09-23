package util

import (
	"regexp"
	"strings"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

func Regexp(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(fl.Param())
	return re.MatchString(fl.Field().String())
}

func StrToTime(layout, str string) time.Time {
	if str == "" {
		return time.Now()
	}
	// *timezone
	loc, _ := time.LoadLocation("Asia/Bangkok")
	t, err := time.ParseInLocation(layout, str, loc)
	if err != nil {
		return time.Now()
	}
	return t
}

func ConcatWithSeperator(values []string, seperator string) string {
	if len(values) == 1 {
		return values[0]
	}
	return strings.Join(values, seperator)
}
