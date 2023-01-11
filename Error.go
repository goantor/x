package x

import "fmt"

var (
	errorsMapping map[IErrno]string
	UNKNOWN_ERRNO FatalErrno = -9999
)

func ErrorsInit(errors map[IErrno]string) {
	if errorsMapping == nil {
		errorsMapping = errors
		errorsMapping[UNKNOWN_ERRNO] = "unknown error"
	}
}

func ErrorSave(no IErrno, message string) {
	errorsMapping[no] = message
}

type Errors map[IErrno]string

type IErrno interface {
	Level() string
}

type Errno int

func (e Errno) Level() string {
	return "base"
}

func (e Errno) Error() string {
	return errorMessage(e)
}

type InfoErrno Errno

func (e InfoErrno) Level() string {
	return "info"
}

func (e InfoErrno) Error() string {
	return errorMessage(e)
}

type TraceErrno Errno

func (e TraceErrno) Level() string {
	return "trace"
}

func (e TraceErrno) Error() string {
	return errorMessage(e)
}

type DebugErrno Errno

func (e DebugErrno) Level() string {
	return "debug"
}

func (e DebugErrno) Error() string {
	return errorMessage(e)
}

type WarnErrno Errno

func (e WarnErrno) Level() string {
	return "warning"
}

func (e WarnErrno) Error() string {
	return errorMessage(e)
}

type ErrorErrno Errno

func (e ErrorErrno) Level() string {
	return "error"
}

func (e ErrorErrno) Error() string {
	return errorMessage(e)
}

type FatalErrno Errno

func (e FatalErrno) Level() string {
	return "fatal"
}

func (e FatalErrno) Error() string {
	return errorMessage(e)
}

type PanicErrno Errno

func (e PanicErrno) Level() string {
	return "panic"
}

func (e PanicErrno) Error() string {
	return errorMessage(e)
}

func GetErrorMessage(no IErrno) string {
	if message, ok := errorsMapping[no]; ok {
		return message
	}

	return "unknown error"
}

type Error struct {
	Code    IErrno // 错误码
	Data    interface{}
	Message string // 内容
}

func New(code IErrno, message string, data interface{}) *Error {
	return &Error{
		Code:    code,
		Data:    data,
		Message: message,
	}
}

func ThrowError(error interface{}, args ...interface{}) {
	var (
		data interface{} = nil
	)

	if len(args) > 0 {
		data = args[0]
	}

	e := getError(error)

	e.Data = data
	log.Auto(e.Code, e.Message, data)
	throw(e)
}

func getError(error interface{}) *Error {
	var (
		e     *Error
		errno IErrno
	)

	switch error.(type) {
	case IErrno:
		errno = error.(IErrno)
		e = &Error{
			Code:    errno,
			Message: errorMessage(errno),
		}
		break
	case Error:
		ie := error.(Error)
		e = &ie
		break
	case *Error:
		e = error.(*Error)
		break
	default:
		e = &Error{
			Code:    UNKNOWN_ERRNO,
			Message: "unknown error",
			Data:    error,
		}
	}

	return e
}

func throw(e *Error) {
	panic(e)
}

func errorMessage(errno IErrno, def ...string) string {
	err, ok := errorsMapping[errno]
	if ok {
		return err
	}

	if len(def) == 0 {
		return errorsMapping[UNKNOWN_ERRNO]
	}

	return def[0]
}

func Recover(fn func(r interface{}, message string)) {
	if r := recover(); r != nil {
		switch r.(type) {
		case IErrno:
			fn(r, errorsMapping[r.(IErrno)])
			break
		case *Error:
			fn(r, "")
			break
		default:
			fn(UNKNOWN_ERRNO, fmt.Sprintf("%+v", r))
		}
	}
}
