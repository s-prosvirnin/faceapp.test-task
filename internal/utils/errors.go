package utils

import "github.com/pkg/errors"

type Typeable interface {
	Type() string
	Error() string
}

type ErrWithType struct {
	err     error
	errType string
}

func NewErrWithType(err error, errType string) ErrWithType {
	return ErrWithType{err: err, errType: errType}
}

func (e ErrWithType) Type() string {
	return e.errType
}

func (e ErrWithType) Error() string {
	return e.err.Error()
}

func (e ErrWithType) Unwrap() error {
	return e.err
}

var ErrEmptyType = NewErrWithType(errors.New(""), "")

func CheckErrorType(err error, errType string) bool {
	for err != nil {
		tErr, ok := err.(Typeable)
		if ok {
			return errType == tErr.Type()
		}
		err = errors.Unwrap(err)
	}

	return false
}

func GetErrorType(err error) Typeable {
	for err != nil {
		tErr, ok := err.(Typeable)
		if ok {
			return tErr
		}
		err = errors.Unwrap(err)
	}

	return ErrEmptyType
}
