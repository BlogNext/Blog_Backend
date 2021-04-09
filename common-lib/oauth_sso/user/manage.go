package user

import (
	"github.com/blog_backend/common-lib/oauth_sso"
	"net/http"
)

//用户资源信息服务
type Manage struct {
	request *request
	//传入参数
	TokenMange *oauth_sso.TokenManage
}

func NewManage(tokenMange *oauth_sso.TokenManage) *Manage {
	return &Manage{
		request:    newRequest(),
		TokenMange: tokenMange,
	}
}

//获取用户信息
func (m *Manage) UserInfo(r *UserInfoResponse) error {
	return m.TokenMange.HttpDoRequest(func(accessToken string) (h *http.Request, entity oauth_sso.DataEntity) {
		return m.request.userInfo(accessToken), r
	})
}
