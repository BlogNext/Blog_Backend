package oauth

import (
	"github.com/blog_backend/common-lib/oauth_sso/core"
	"net/http"
	"net/url"
	"strings"
)

//协议
type request struct {
	url *core.OauthSSOUrl
}

func newRequest() *request {
	return &request{
		url: core.NewUrl(),
	}
}

//创建预授权码
func (r *request) createPreAuthCode(nickname, password, clientId, redirectUrl string) *http.Request {
	values := url.Values{}
	values.Set("nickname", nickname)
	values.Set("password", password)
	values.Set("client_id", clientId)
	values.Set("redirect_url", redirectUrl)
	req, _ := http.NewRequest(http.MethodPost, r.url.GetUrl("api/oauth/create_pre_auth_code"), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

//授权码换取accessToken
func (r *request) preAuthCodeAccessToken(preAuthCode, clientId, clientSecret string) *http.Request {
	values := url.Values{}
	values.Set("pre_auth_code", preAuthCode)
	values.Set("client_id", clientId)
	values.Set("client_secret", clientSecret)
	req, _ := http.NewRequest(http.MethodPost, r.url.GetUrl("api/oauth/pre_auth_code_access_token"), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

//授权码换取accessToken
func (r *request) refreshToken(refreshToken string) *http.Request {
	values := url.Values{}
	values.Set("refresh_token", refreshToken)
	req, _ := http.NewRequest(http.MethodPost, r.url.GetUrl("api/oauth/refresh_token"), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
