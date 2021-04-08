package oauth_sso

import (
	"github.com/blog_backend/common-lib/oauth_sso/oauth"
	"github.com/blog_backend/common-lib/oauth_sso/user"
)

//入口函数

//oauthsso的服务器配置信息
var scheme string
var host string

func init() {
	if scheme == "" {
		scheme = "http"
		host = "127.0.0.1:8084"
	}
}

//获取OAuth的Manage
func GetAuthManage() *oauth.Manage {
	return &oauth.Manage{Request: oauth.Request{
		Url: NewUrl(scheme, host),
	}}
}

//获取User的Manage
func GetUserManageByPreAuthCode(preAuthCode, clientId, clientSecret string) *user.Manage {
	
	return nil
}
