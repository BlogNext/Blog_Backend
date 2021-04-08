package oauth

//创建预授权码的响应
type CreatePreAuthCodeResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
