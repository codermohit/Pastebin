package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid() return true if FieldErrors map doesn't contain any errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= 8
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// add a key and error message to the map
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField() adds an error message to the FieldErrors map only if a validation check is not 'ok'
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

//Create a AddNonFieldError() helper for adding error messages to the new NonFieldErrors slice 
func (v *Validator) AddNonFieldError(message string) {
  v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// NotBlank() checks for an empty string, true for not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if the value contains no more than n chars
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValues() returns true if a value is in a list of specific permitted values
func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
