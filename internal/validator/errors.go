package validator

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Code    string
	Message string
}

type ValidationErrors struct {
	Errors []ValidationError
}

func (e *ValidationErrors) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d validation errors:\n\n", len(e.Errors)))
	for _, err := range e.Errors {
		sb.WriteString(fmt.Sprintf("[%s]\n%s\n\n", err.Code, err.Message))
	}
	return sb.String()
}
