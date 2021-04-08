package user

import (
	"github.com/blog_backend/common-lib/oauth_sso"
	"net/http"
)

type Manage struct {
	request request
	//授权的用户
	AccessToken oauth_sso.AccessToken
}

//创建预授权码
func (m *Manage) UserInfo(r *UserInfoResponse) oauth_sso.RequestInitFunc {
	return func() (*http.Request, oauth_sso.DataEntity) {
		return m.request.userInfo(m.AccessToken), r
	}
}
