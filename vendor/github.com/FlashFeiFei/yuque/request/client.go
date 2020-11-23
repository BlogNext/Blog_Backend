package request

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	BASE_URI = "https://www.yuque.com/api/v2"
)

//选项
type Option func(request *http.Request)

//post、put请求的时候设置Content-Type 是 application/json 类型
func contentTypeOp() Option {

	return func(request *http.Request) {
		request.Header.Add("Content-Type", "application/json")
	}

}

//这是token
func NewTokenOp(token string) Option {
	return func(request *http.Request) {
		request.Header.Add("X-Auth-Token", token)
	}
}

//设置用户代理
func NewUserAgentOp(user_agent string) Option {
	return func(request *http.Request) {

		request.Header.Add("User-Agent", user_agent)
	}
}

//请求客户端
type ClientHandle struct {
	Op      []Option      //选项
	BaseUri string        //baseuri
	Uri string            //uri
	Method string        //请求方法
	Token   string        //用户token
	Body io.Reader       //请求体
}

func NewClientHandle(token,uri ,method string, op []Option) (c*ClientHandle){
	c = new(ClientHandle)
	c.Op = op
	c.BaseUri = BASE_URI
	c.Uri = uri
	c.Method = method
	c.Token = token
	return c
}

func (c *ClientHandle) Request(response interface{})(err error) {
	client := http.Client{}

	req ,err :=http.NewRequest(c.Method,fmt.Sprintf("%s%s",c.BaseUri,c.Uri),c.Body)
	if err != nil {
		return
	}

	//请求前做一些处理

	if c.Op == nil {
		c.Op = make([]Option,0)
	}
	//所有api都需要带上token
	c.Op = append(c.Op,NewTokenOp(c.Token))

	//请求方法添加contentType
	switch req.Method {
	case http.MethodPut:
		c.Op = append(c.Op, contentTypeOp())
	case http.MethodPost:
		c.Op = append(c.Op, contentTypeOp())
	}


	//回调请求前回调下选项设计
	for _, op := range c.Op {
		op(req)
	}

	//请求
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	//读取数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	//解码
	err = json.Unmarshal(body,response)

	return err
}
