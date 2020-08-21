package exception

const (
	//框架错误码
	VALIDATE_ERR = iota + 1 //验证错误
	//token错误码
	TOKEN_EXPIRED = iota + 1000
	TOKEN_ILLEGALITY
)
