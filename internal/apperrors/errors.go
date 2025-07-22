package apperrors

import (
	"fmt"
	"github.com/google/uuid"
)

type NotFoundError struct {
	Resource string
	ID       uuid.UUID
}

type DuplicateError struct {
	Field string
}

type InvalidInputError struct {
	Message string
}

func (ie *InvalidInputError) Error() string {
	return ie.Message
}

func (de *DuplicateError) Error() string {
	return fmt.Sprintf("%s already exists", de.Field)
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", n.Resource, n.ID.String())
}

func NewNotFoundError(resource string, id uuid.UUID) error {
	return &NotFoundError{Resource: resource, ID: id}
}

func NewDuplicateError(field string) error {
	return &DuplicateError{Field: field}
}

func NewInvalidInputError(msg string) error {
	return &InvalidInputError{Message: msg}
}
