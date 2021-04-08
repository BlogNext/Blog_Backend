package oauth

import (
	"github.com/blog_backend/common-lib/oauth_sso"
	"net/http"
)

type Manage struct {
	Request Request
}

//创建预授权码
func (m *Manage) CreatePreAuthCode(nickname, password, clientId, redirectUrl string, r *CreatePreAuthCodeResponse) oauth_sso.RequestInitFunc {
	return func() (*http.Request, oauth_sso.DataEntity) {
		return m.Request.createPreAuthCode(nickname, password, clientId, redirectUrl), r
	}
}

//预授权码换取token
func (m *Manage) PreAuthCodeAccessToken(preAuthCode, clientId, clientSecret string, r *PreAuthCodeAccessTokenResponse) oauth_sso.RequestInitFunc {
	return func() (*http.Request, oauth_sso.DataEntity) {
		return m.Request.preAuthCodeAccessToken(preAuthCode, clientId, clientSecret), r
	}
}

//refreshToken刷新
func (m *Manage) RefreshToken(refreshToken string, r *RefreshTokenResponse) oauth_sso.RequestInitFunc {
	return func() (*http.Request, oauth_sso.DataEntity) {
		return m.Request.refreshToken(refreshToken), r
	}
}
