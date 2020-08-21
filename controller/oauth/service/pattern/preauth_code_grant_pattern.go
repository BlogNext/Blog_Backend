package pattern

import (
	"github.com/blog_backend/controller/oauth/service"
)

const (
	PREAUTH_CODE_STEP = iota //预授权码步骤
	TOKEN_STEP               //预授权码换取token步骤
)

//授权码模式
type PreauthCodeGrantPattern struct {
	*BaseGrantPattern
	step         int
	preauth_code string //预授权码
}

func NewPreauthCodeGrantPattern(step int, preauth_code string, client_and_user *service.ClientAndUser) *PreauthCodeGrantPattern {

	switch step {
	case PREAUTH_CODE_STEP:
	case TOKEN_STEP:
		break
	default:
		panic("PreauthCodeGrantPattern未支持的步骤")
	}

	return &PreauthCodeGrantPattern{
		BaseGrantPattern: NewBaseGrantPattern(client_and_user),
		step:             step,
		preauth_code:     preauth_code,
	}
}

func (p *PreauthCodeGrantPattern) Oauth(account_is_valid ...service.ValidateAccountIsValid) (*service.Token, error) {
	switch p.step {
	case PREAUTH_CODE_STEP:
		preauth_code, err := service.GetPreauthCode(p.BaseGrantPattern.ClientAndUser, account_is_valid...)
		if err != nil {
			return nil, err
		}
		return preauth_code, nil
	case TOKEN_STEP:
		token, err := service.GetAccessToken(p.preauth_code, p.BaseGrantPattern.ClientAndUser, account_is_valid...)
		if err != nil {
			return nil, err
		}
		return token, nil
	}

	panic("检查代码！")
}
