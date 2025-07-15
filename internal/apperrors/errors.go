package apperrors

import (
	"fmt"
	"strings"
)

type NotFoundError struct {
	Resource string
	ID       int
}

type ValidationError struct {
	Messages []string
}

func (v ValidationError) Error() string {
	return "invalid input:" + strings.Join(v.Messages, ", ")
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", n.Resource, n.ID)
}

func NewNotFoundError(resource string, id int) error {
	return &NotFoundError{Resource: resource, ID: id}
}

func NewValidationError(messages ...string) error {
	return &ValidationError{Messages: messages}
}
