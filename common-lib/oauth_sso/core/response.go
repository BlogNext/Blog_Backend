package core

//数据
type DataEntity interface {
	GetData() interface{}
}

//响应
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (r *Response) SetData(dataEntity DataEntity) {
	r.Data = dataEntity.GetData()
}
