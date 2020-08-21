package pattern

import (
	"github.com/blog_backend/controller/oauth/service"
)

//应用授信模式
type ClientCredentialGrantPattern struct {
	BaseGrantPattern
}

func (c *ClientCredentialGrantPattern) Oauth(account_is_valid ...service.ValidateAccountIsValid) (*service.Token, error) {
	return nil, nil
}
