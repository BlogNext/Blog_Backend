package oauth_sso

import (
	"fmt"
	"net/url"
)

type Url url.URL

//url创建
func NewUrl(scheme, host string) *Url {
	return &Url{
		Scheme: scheme,
		Host:   host,
	}
}

//获取url配置
func (u *Url) GetUrl(uri string) string {
	return fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, uri)
}
