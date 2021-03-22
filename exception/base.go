package exception

type MyException interface {
	error
	SetErrorCode(errorCode int)
	SetErrorMsg(errMsg string)
	GetErrorCode() int
}

type BaseException struct {
	errorMsg  string
	errorCode int
}

func NewException(errorCode int, errorMsg string) *BaseException {
	return &BaseException{
		errorMsg:  errorMsg,
		errorCode: errorCode,
	}
}

func (b *BaseException) SetErrorCode(errorCode int) {
	b.errorCode = errorCode
}

func (b *BaseException) SetErrorMsg(errorMsg string) {
	b.errorMsg = errorMsg
}

func (b *BaseException) GetErrorCode() int {
	return b.errorCode
}

func (b *BaseException) Error() string {
	return b.errorMsg
}
