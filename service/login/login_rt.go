package login

import (
	"errors"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity/login/front"
	"github.com/blog_backend/model"
	"gorm.io/gorm"
	"strings"
)

//前端登录
type LoginRtService struct {
}

//语雀登录
//login语雀的login
//password登录密码
func (u *LoginRtService) LoginByYuque(login, password string) (login_entity *front.LoginEntity) {
	db := mysql.GetDefaultDBConnect()
	model := new(model.UserYuQueModel)
	query_result := db.Where("login = ?", login).First(model)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		panic("为找到用户login:" + login)
	}

	if strings.Compare(password, "xiaochen123") != 0 {
		panic("密码不正确")
	}

	
	return login_entity
}
