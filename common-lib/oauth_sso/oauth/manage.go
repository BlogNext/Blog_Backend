package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blog_backend/common-lib/oauth_sso"
	"io/ioutil"
	"net/http"
)

type ManageRequestInitFunc func() (*http.Request, oauth_sso.DataEntity)

//授权的服务
type Manage struct {
	//private的，不暴露的参数
	request *request
	//必填
	clientId     string
	clientSecret string
}

//创建manage
func NewManage(clientId, clientSecret string) *Manage {
	return &Manage{
		request:      newRequest(),
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

//创建预授权码
func (m *Manage) CreatePreAuthCode(nickname, password, redirectUrl string, r *CreatePreAuthCodeResponse) error {
	return m.HttpDoRequest(func() (*http.Request, oauth_sso.DataEntity) {
		return m.request.createPreAuthCode(nickname, password, m.clientId, redirectUrl), r
	})
}

//预授权码换取token
func (m *Manage) PreAuthCodeAccessToken(preAuthCode string, r *PreAuthCodeAccessTokenResponse) error {
	return m.HttpDoRequest(func() (*http.Request, oauth_sso.DataEntity) {
		return m.request.preAuthCodeAccessToken(preAuthCode, m.clientId, m.clientSecret), r
	})
}

//refreshToken刷新
func (m *Manage) RefreshToken(refreshToken string, r *RefreshTokenResponse) error {
	return m.HttpDoRequest(func() (*http.Request, oauth_sso.DataEntity) {
		return m.request.refreshToken(refreshToken), r
	})
}

//http请求
func (m *Manage) HttpDoRequest(manageRequestInitFunc ManageRequestInitFunc) (err error) {

	defer func() {
		if myException := recover(); myException != nil {
			err = myException.(error)
			return
		}
	}()

	//调用回调得到request和响应
	request, dataEntity := manageRequestInitFunc()
	//发送请求
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	//响应
	entityResponse := new(oauth_sso.Response)
	entityResponse.SetData(dataEntity)
	err = json.Unmarshal(data, entityResponse)
	if err != nil {
		//序列化失败
		return err
	}

	if entityResponse.Code != 0 {
		return errors.New(fmt.Sprintf("code:%d,错误信息:%s", entityResponse.Code, entityResponse.Msg))
	}

	return nil
}
