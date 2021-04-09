package user

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
func (r *request) userInfo(accessToken string) *http.Request {
	values := url.Values{}
	values.Set("access_token", accessToken)
	req, _ := http.NewRequest(http.MethodPost, r.url.GetUrl("api/resource/user/user_info"), strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
