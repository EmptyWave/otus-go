package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrValidateWrong = errors.New("validate wrong")
	ErrValidate      = errors.New("validated error")
	ErrNoStruct      = errors.New("invalid data type")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errStr := strings.Builder{}
	for _, err := range v {
		if errors.Is(err.Err, ErrValidate) {
			errStr.WriteString(fmt.Sprintf("Field: %s: %s\n", err.Field, err.Err))
		}
	}
	return errStr.String()
}
