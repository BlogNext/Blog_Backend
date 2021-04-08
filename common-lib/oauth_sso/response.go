package oauth_sso

//数据
type DataEntity interface{}

//响应
type Response struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data DataEntity `json:"data,omitempty"`
}

func (r *Response) SetData(dataEntity DataEntity) {
	r.Data = dataEntity
}
