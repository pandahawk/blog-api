package apperrors

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type NotFoundError struct {
	Resource string
	ID       uuid.UUID
}

type ValidationError struct {
	Messages []string
}

func (v ValidationError) Error() string {
	return "invalid input: " + strings.Join(v.Messages, ", ")
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", n.Resource, n.ID.String())
}

func NewNotFoundError(resource string, id uuid.UUID) error {
	return &NotFoundError{Resource: resource, ID: id}
}

func NewValidationError(messages ...string) error {
	return &ValidationError{Messages: messages}
}
