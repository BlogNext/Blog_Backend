package oauth_sso

import (
	"fmt"
	"net/url"
)

//oauthSSO的服务器配置信息
var oauthSSOConfig map[string]string

func SetOauthSSoSchemeConfig(scheme string) {
	oauthSSOConfig["scheme"] = scheme
}

func SetOauthSSoHostConfig(host string) {
	oauthSSOConfig["host"] = host
}

func init() {

	if oauthSSOConfig == nil {
		oauthSSOConfig = make(map[string]string)
		//协议
		SetOauthSSoSchemeConfig("")
		//地址
		SetOauthSSoHostConfig("")
	}
}

//oauthSSo配置
type OauthSSOUrl url.URL

//url创建
func NewUrl() *OauthSSOUrl {
	return &OauthSSOUrl{
		Scheme: oauthSSOConfig["scheme"],
		Host:   oauthSSOConfig["host"],
	}
}

//获取url配置
func (o *OauthSSOUrl) GetUrl(uri string) string {
	return fmt.Sprintf("%s://%s/%s", o.Scheme, o.Host, uri)
}
