package exception

type MyException interface {
	error
	SetErrorCode(error_code int)
	SetErrorMsg(err_msg string)
	GetErrorCode() int
}

type BaseException struct {
	error_msg  string
	error_code int
}

func NewException(error_code int, error_msg string) *BaseException {
	return &BaseException{
		error_msg:  error_msg,
		error_code: error_code,
	}
}

func (b *BaseException) SetErrorCode(error_code int) {
	b.error_code = error_code
}

func (b *BaseException) SetErrorMsg(err_msg string) {
	b.error_msg = err_msg
}

func (b *BaseException) GetErrorCode() int {
	return b.error_code
}

func (b *BaseException) Error() string {
	return b.error_msg
}
