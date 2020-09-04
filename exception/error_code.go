package exception

const (
	//框架错误码
	VALIDATE_ERR = iota + 1 //验证错误
	//token错误码
	TOKEN_EXPIRED = iota + 10000
	TOKEN_ILLEGALITY
	//数据错误码
	DATA_BASE_ERROR_EXEC = iota + 20000  //数据库错误的执行
	//es错误码
	ES_ERROR_CONNET = iota + 30000   //es链接错误
)
