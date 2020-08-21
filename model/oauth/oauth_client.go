package oauth

import (
	"github.com/blog_backend/model"
	"time"
)

type OauthClient struct {
	model.BaseModel
	ClientName      string     `gorm:"cloumn:client_name"`
	ClientAppId     string     `gorm:"cloumn:client_app_id"`
	ClientAppSecret string     `gorm:"cloumn:client_app_secret"`
	RedirectUrl     string     `gorm:"cloumn:redirect_url"`
	Year            int        `gorm:"cloumn:year"`
	Month           time.Month `gorm:"cloumn:month"`
}

func (OauthClient) TableName() string {
	return "oauth_client"
}
