package oauth

//创建预授权码的响应
type CreatePreAuthCodeResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (c *CreatePreAuthCodeResponse) GetData() interface{} {
	return c
}

//预授权码换取token
type PreAuthCodeAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (p *PreAuthCodeAccessTokenResponse) GetData() interface{} {
	return p
}

//refreshToken刷新
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshTokenResponse) GetData() interface{} {
	return r
}
