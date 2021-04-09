package user

import (
	"github.com/blog_backend/common-lib/oauth_sso/core"
	"github.com/blog_backend/common-lib/oauth_sso/token"
	"net/http"
)

//用户资源信息服务
type Manage struct {
	request *request
	//传入参数
	TokenMange *token.TokenManage
}

func NewManage(tokenMange *token.TokenManage) *Manage {
	return &Manage{
		request:    newRequest(),
		TokenMange: tokenMange,
	}
}

//获取用户信息
func (m *Manage) UserInfo(r *UserInfoResponse) error {
	return m.TokenMange.HttpDoRequest(func(accessToken string) (h *http.Request, entity core.DataEntity) {
		return m.request.userInfo(accessToken), r
	})
}
