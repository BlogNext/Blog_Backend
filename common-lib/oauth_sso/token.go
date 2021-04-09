package oauth_sso

import (
	"encoding/json"
	"fmt"
	"github.com/blog_backend/common-lib/oauth_sso/oauth"
	"github.com/blog_backend/exception"
	"io/ioutil"
	"net/http"
)

type RequestInitFunc func(accessToken string) (*http.Request, DataEntity)

//用户的token管理
type TokenManage struct {
	//必填的
	refreshToken string
	oauthManage  *oauth.Manage
	//非必填的
	accessToken string
}

//创建一个token
func NewTokenManage(refreshToken, clientId, clientSecret string) *TokenManage {
	return &TokenManage{
		refreshToken: refreshToken,
		oauthManage:  oauth.NewManage(clientId, clientSecret),
	}
}

//获取accessToken
//是否通过isRefreshToken强制刷新accessToken,true 强刷,fasle强刷
func (m *TokenManage) GetAccessToken(isRefreshToken bool) string {

	if isRefreshToken == true {
		//强刷
		refreshToken := new(oauth.RefreshTokenResponse)
		err := m.oauthManage.RefreshToken(m.refreshToken, refreshToken)
		if err != nil {
			panic(err)
		}
		m.accessToken = refreshToken.AccessToken
	}

	return m.accessToken
}

//需要依赖accessToken的请求有这里发送
func (m *TokenManage) HttpDoRequest(requestInitFunc RequestInitFunc) error {

	err := m.httpRequest(requestInitFunc)
	if err != nil {
		//刷新accessToken再试一次
		m.GetAccessToken(true)
		err = m.httpRequest(requestInitFunc)
	}

	if err != nil {
		return err
	}

	return nil
}

//http请求
//isResult true成功， false不成功
func (m *TokenManage) httpRequest(requestInitFunc RequestInitFunc) (err error) {

	defer func() {
		if errPanic := recover(); errPanic != nil {
			err = errPanic.(error)
			return
		}
	}()

	//调用回调得到request和响应
	request, dataEntity := requestInitFunc(m.accessToken)
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
	entityResponse := new(Response)
	entityResponse.SetData(dataEntity)
	err = json.Unmarshal(data, entityResponse)
	if err != nil {
		//序列化失败
		panic(err)
	}

	if entityResponse.Code != 0 {
		panic(exception.NewException(exception.VALIDATE_ERR, fmt.Sprintf("code:%d,错误信息:%s", entityResponse.Code, entityResponse.Msg)))
	}

	return nil
}
