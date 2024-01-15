package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Valid checks if the Validator object is valid.
//
// It returns a boolean indicating whether there are any field errors.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError adds an error message for a specific field key in the Validator.
//
// Parameters:
// - key: the key of the field to add the error message for (string).
// - message: the error message to be added (string).
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField adds an error to the Validator for the given field if the condition is not met.
//
// Parameters:
// - ok (bool): The condition to check.
// - key (string): The key of the field.
// - message (string): The error message to associate with the field.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank checks if a string value is not empty or only whitespace.
//
// value: the string value to be checked.
// Returns true if the string is not empty or only whitespace, false otherwise.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars checks if the length of the given string is less than or equal to a specified number.
//
// Parameters:
// - value: the string to be checked.
// - n: the maximum number of characters allowed.
//
// Returns:
// - bool: true if the length of the string is less than or equal to n, false otherwise.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// PermittedInt checks if a given value is present in a list of permitted values.
//
// Parameters:
// - value: the value to be checked.
// - permittedValues: the list of permitted values.
//
// Returns:
// - bool: true if the value is present in the list, false otherwise.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
