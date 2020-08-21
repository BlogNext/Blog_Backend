package pattern

import (
	"github.com/blog_backend/controller/oauth/service"
)

//简化模式
type ImplicitGrantPattern struct {
	*BaseGrantPattern
}

func NewImplicitGrantPattern(client_and_user *service.ClientAndUser) *ImplicitGrantPattern {

	return &ImplicitGrantPattern{
		BaseGrantPattern: NewBaseGrantPattern(client_and_user),
	}
}

func (i *ImplicitGrantPattern) Oauth(account_is_valid ...service.ValidateAccountIsValid) (*service.Token, error) {
	preauth_code, err := service.GetPreauthCode(i.BaseGrantPattern.ClientAndUser, account_is_valid...)
	if err != nil {
		return nil, err
	}

	token, err := service.GetAccessToken(preauth_code.PreauthCode, i.BaseGrantPattern.ClientAndUser, account_is_valid...)
	if err != nil {
		return nil, err
	}

	return token, nil
}
