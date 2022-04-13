package errs

type ResponseError struct {
	Message string `json:"message"`
}

func (e *ResponseError) Error() string {
	return e.Message
}

// EmptyBody is an empty instance of ResponseError type, we need it for the cases in which empty body
// of the request is allowed
type EmptyBody ResponseError

func (e *EmptyBody) Error() string {
	return e.Message
}
