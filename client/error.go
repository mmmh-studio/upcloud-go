package client

import (
	"fmt"
	"reflect"
)

type Error struct {
	code    int
	message string
	body    []byte
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.message)
}

func (e *Error) Body() string {
	return string(e.body)
}

type FieldError struct {
	// Name of the field.
	Name string
	// Descirption of the error.
	Description string
}

type ValidationError struct {
	// Name of the payload the validation failed on.
	Name string
	// FieldErrors encountered during validation.
	FieldErrors []FieldError
}

func (e *ValidationError) Error() string {
	str := fmt.Sprintf("ValidationError: %s", e.Name)

	for _, err := range e.FieldErrors {
		str = fmt.Sprintf("%s\n\t%s - %s", str, err.Name, err.Description)
	}

	return str
}

func typeName(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
