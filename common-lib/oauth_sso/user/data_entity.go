package user

//创建预授权码的响应
type UserInfoResponse struct {
	Id       uint64 `json:"id"`
	Nickname string `json:"nickname"`
}

func (u *UserInfoResponse) GetData() interface{} {
	return u
}
