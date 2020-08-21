package pattern

import (
	"errors"
	"github.com/blog_backend/controller/oauth/service"
)

type Oauth interface {
	Oauth(account_is_valid ...service.ValidateAccountIsValid) (token *service.Token, err error)
}

type BaseGrantPattern struct {
	*service.ClientAndUser
}

func NewBaseGrantPattern(client_and_user *service.ClientAndUser) *BaseGrantPattern {
	return &BaseGrantPattern{client_and_user}
}

func (b *BaseGrantPattern) Oauth(account_is_valid ...service.ValidateAccountIsValid) (*service.Token, error) {
	return nil, errors.New("请重写具体的授权模式")
}
