package oauth_sso

import (
	"encoding/json"
	"fmt"
	"github.com/blog_backend/exception"
	"io/ioutil"
	"net/http"
)

type RequestInitFunc func() (*http.Request, *Response)

//发送请求
func HttpDoRequest(requestInitFunc RequestInitFunc) {
	//调用回调得到request和响应
	request, entityResponse := requestInitFunc()

	//发送请求
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(exception.NewException(exception.VALIDATE_ERR, err.Error()))
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	//响应
	err = json.Unmarshal(data, entityResponse)
	if err != nil {
		//序列化失败
		panic(err)
	}

	if entityResponse.Code != 0 {
		panic(exception.NewException(exception.VALIDATE_ERR, fmt.Sprintf("code:%d,错误信息:%s", entityResponse.Code, entityResponse.Msg)))
	}
}
