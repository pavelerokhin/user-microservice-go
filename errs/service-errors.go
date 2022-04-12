package errs

type ResponseError struct {
	Message string `json:"message"`
}

func (e *ResponseError) Error() string {
	return e.Message
}

type EmptyBody ResponseError

func (e *EmptyBody) Error() string {
	return e.Message
}
