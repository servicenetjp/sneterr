// Package odinerr represents API error interface
package odinerr

import (
	"fmt"
	"path"
	"runtime"
)

// An Error wraps lower level errors with code, message and an original error.
type Error interface {
	// Satisfy the generic error interface.
	error

	// Returns the short phrase depicting the classification of the error.
	Code() string

	// Returns the error details message.
	Message() string

	// Returns the original error if one was set.  Nil is returned if not set.
	OrigErr() error
}

// SprintError returns a string of the formatted error code.
//
// Both extra and origErr are optional.  If they are included their lines
// will be added, but if they are not included their lines will be ignored.
func SprintError(code, message, extra string, origErr error) string {
	msg := fmt.Sprintf("%s: %s", code, message)
	if extra != "" {
		msg = fmt.Sprintf("%s\n\t%s", msg, extra)
	}
	if origErr != nil {
		msg = fmt.Sprintf("%s\ncaused by: %s", msg, origErr.Error())
	}
	return msg
}

// A baseError wraps the code and message which defines an error. It also
// can be used to wrap an original error object.
//
// Should be used as the root for errors satisfying the odinerr.Error. Also
// for any error which does not fit into a specific error wrapper type.
type baseError struct {
	// Classification of error
	code string

	// Detailed information about error
	message string

	// Optional original error. O que causou o erro
	err error

	//arquivo onde foi criado o erro
	file string

	//linha do arquivo onde foi criado
	linhaArquivo int
}

// newBaseError returns an error object for the code, message, and errors.
//
// code is a short no whitespace phrase depicting the classification of
// the error that is being created.
//
// message is the free flow string containing detailed information about the
// error.
func newBaseError(code, message string, origErr error, file string, linha int) *baseError {
	b := &baseError{
		code:    code,
		message: message,
		err:     origErr,
	}

	return b
}

// Error returns the string representation of the error.
//
// Satisfies the error interface.
func (b baseError) Error() string {
	var causaErro string
	if b.err != nil {
		causaErro = b.err.Error()
	}
	return fmt.Sprintf("(%s:%d) (code:%s) (msg:%s) (err:%s)",
		b.file, b.linhaArquivo, b.code, b.message, causaErro)
}

// String returns the string representation of the error.
// Alias for Error to satisfy the stringer interface.
func (b baseError) String() string {
	return b.Error()
}

// Code returns the short phrase depicting the classification of the error.
func (b baseError) Code() string {
	return b.code
}

// Message returns the error details message.
func (b baseError) Message() string {
	return b.message
}

//OrigErr ...
func (b baseError) OrigErr() error {
	return b.err
}

// New returns an Error object described by the code, message, and origErr.
//
// If origErr satisfies the Error interface it will not be wrapped within a new
// Error object and will instead be returned.
func New(code, message string, origErr error) Error {
	_, file, line, _ := runtime.Caller(1)
	_, nomeArquivo := path.Split(file)

	return newBaseError(code, message, origErr, nomeArquivo, line)
}
