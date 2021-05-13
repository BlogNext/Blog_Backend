package front


//oauthSSo单点登录请求
type BlogNextPreCodeRequest struct {
	PreCode string `form:"pre_code" json:"pre_code" binding:"required"`
}
