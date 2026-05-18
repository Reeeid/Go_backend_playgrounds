package post

import "strconv"

type LengthError struct {
	FieldName string
	Min       int
	Max       int
}

type NilError struct {
	FieldName string
}

type StatusError struct {
	FieldName string
}

func (e *LengthError) Error() string {
	return e.FieldName + " must be between " + strconv.Itoa(e.Min) + " and " + strconv.Itoa(e.Max) + " characters"
}

func (e *NilError) Error() string {
	return e.FieldName + " cannot be nil"
}

func (e *StatusError) Error() string {
	return e.FieldName + " must be one of: " + string(StatusDraft) + ", " + string(StatusPublished) + ", " + string(StatusArchived)
}
