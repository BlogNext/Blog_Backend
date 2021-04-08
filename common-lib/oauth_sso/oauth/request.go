package oauth

import (
	"github.com/blog_backend/common-lib/oauth_sso"
	"net/http"
	"net/url"
	"strings"
)

//协议
type request struct {
	url *oauth_sso.Url
}

//创建预授权码
func (o *request) CreatePreAuthCode(nickname, password, clientId, redirectUrl string) *http.Request {
	values := url.Values{}
	values.Set("nickname", nickname)
	values.Set("password", password)
	values.Set("client_id", clientId)
	values.Set("redirect_url", redirectUrl)
	req, _ := http.NewRequest(http.MethodPost, o.url.GetUrl("api/oauth/create_pre_auth_code"), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

//授权码换取accessToken
func (o *request) preAuthCodeAccessToken(preAuthCode, clientId, clientSecret string) *http.Request {
	values := url.Values{}
	values.Set("pre_auth_code", preAuthCode)
	values.Set("client_id", clientId)
	values.Set("client_secret", clientSecret)
	req, _ := http.NewRequest(http.MethodPost, o.url.GetUrl("api/oauth/pre_auth_code_access_token"), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

//授权码换取accessToken
func (o *request) RefreshToken(refreshToken string) *http.Request {
	values := url.Values{}
	values.Set("refresh_token", refreshToken)
	req, _ := http.NewRequest(http.MethodPost, o.url.GetUrl("api/oauth/refresh_token"), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func (o *request) VerifyAccessToken(accessToken string) *http.Request {
	values := url.Values{}
	values.Set("access_token", accessToken)
	req, _ := http.NewRequest(http.MethodPost, o.url.GetUrl("api/oauth/verify_access_token"), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
