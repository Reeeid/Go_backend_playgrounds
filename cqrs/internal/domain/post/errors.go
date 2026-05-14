package post

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
	return e.FieldName + " must be between " + string(e.Min) + " and " + string(e.Max) + " characters"
}

func (e *NilError) Error() string {
	return e.FieldName + " cannot be nil"
}

func (e *StatusError) Error() string {
	return e.FieldName + " must be one of: " + string(StatusDraft) + ", " + string(StatusPublished) + ", " + string(StatusArchived)
}
