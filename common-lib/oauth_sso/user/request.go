package user

import (
	"github.com/blog_backend/common-lib/oauth_sso"
	"net/http"
	"net/url"
	"strings"
)

//协议
type request struct {
	url *oauth_sso.OauthSSOUrl
}

func newRequest() *request {
	return &request{
		url: oauth_sso.NewUrl(),
	}
}

//创建预授权码
func (r *request) userInfo(accessToken string) *http.Request {
	values := url.Values{}
	values.Set("refresh_token", accessToken)
	req, _ := http.NewRequest(http.MethodPost, r.url.GetUrl("api/resource/user/user_info"), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
