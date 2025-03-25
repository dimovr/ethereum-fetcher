package errors

type ServiceError struct {
	Msg string
}

func NewServiceError(msg string) ServiceError {
	return ServiceError{Msg: msg}
}

func (e ServiceError) Error() string {
	return e.Msg
}
