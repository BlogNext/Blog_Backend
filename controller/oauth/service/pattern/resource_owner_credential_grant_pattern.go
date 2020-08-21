package pattern

import (
	"github.com/blog_backend/controller/oauth/service"
)

//用户授信模式
type ResourceOwnerCredentialGrantPattern struct {
	BaseGrantPattern
}

func (r *ResourceOwnerCredentialGrantPattern) Oauth(account_is_valid ...service.ValidateAccountIsValid) (*service.Token, error) {
	return nil, nil
}
