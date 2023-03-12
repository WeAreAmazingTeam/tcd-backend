package helper

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LcFirst(str string) string {
	return strings.ToLower(str)
}

func Split(src string) string {
	if !utf8.ValidString(src) {
		return src
	}

	var (
		entries   []string
		runes     [][]rune
		lastClass int = 0
		class     int = 0
	)

	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}

		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}

		lastClass = class
	}

	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}

	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}

	for index, word := range entries {
		if index == 0 {
			entries[index] = UcFirst(word)
		} else {
			entries[index] = LcFirst(word)
		}
	}

	justString := strings.Join(entries, " ")

	return justString
}

func ValidationErrorToText(e validator.FieldError) string {
	word := Split(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required!", word)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s!", word, e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s!", word, e.Param())
	case "email":
		return "invalid email format!"
	case "len":
		return fmt.Sprintf("%s must be %s characters long!", word, e.Param())
	}

	return fmt.Sprintf("%s is not valid!", word)
}

func FormatValidationError(err error) (errors []string) {
	errType := fmt.Sprintf("%T", err)

	if errType == "validator.ValidationErrors" {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationErrorToText(e))
		}
		return errors
	}

	errors = append(errors, err.Error())
	return errors
}
