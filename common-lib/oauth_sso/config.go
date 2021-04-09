package oauth_sso

import (
	"fmt"
	"github.com/blog_backend/common-lib/config"
	"log"
	"net/url"
)

//oauthSSO的服务器配置信息
var oauthSSOConfig map[string]string

func init() {

	//导入配置
	config.LoadConfig("server", "config", "yaml")
	serverConfig, err := config.GetConfig("server")
	if err != nil {
		log.Println(err)
	}
	oauthSSOInfo := serverConfig.GetStringMap("oauthSSO")

	if oauthSSOConfig == nil {
		oauthSSOConfig = make(map[string]string)
		//协议
		oauthSSOConfig["scheme"] = oauthSSOInfo["scheme"].(string)
		//地址
		oauthSSOConfig["host"] = oauthSSOInfo["host"].(string)
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
