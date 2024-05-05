package models

type ResponseError string

func (r ResponseError) Error() string {
	return string(r)
}

const (
	ErrBadRequest   = ResponseError("Bad Request")
	ErrTokenInvalid = ResponseError("Token Invalid")
	ErrorInternal   = ResponseError("Internal Server Error")
)
