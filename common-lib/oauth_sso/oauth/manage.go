package oauth

import (
	"github.com/blog_backend/common-lib/oauth_sso"
	"net/http"
)

type Manage struct {
	request request
}

//创建预授权码
func (m *Manage) CreatePreAuthCode(nickname, password, clientId, redirectUrl string, r *CreatePreAuthCodeResponse) oauth_sso.RequestInitFunc {
	return func() (*http.Request, *oauth_sso.Response) {
		response := new(oauth_sso.Response)
		response.SetData(r)
		return m.request.CreatePreAuthCode(nickname, password, clientId, redirectUrl), response
	}
}
