package user

import (
	"errors"
	"fmt"
	"github.com/FlashFeiFei/yuque/response"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/model"
	"gorm.io/gorm"
)

/**
用户后台服务
*/
type UserBkService struct {
}

func (s *UserBkService) UpdateUserByYuqueWebHook(user *response.UserSerializer) *model.UsereModel {
	db := mysql.GetDefaultDBConnect()
	user_model := new(model.UsereModel)

	//更新用户
	err := db.Transaction(func(tx *gorm.DB) error {

		//找到语雀用户
		user_yuque_model := new(model.UsereYuQueModel)
		query_result := tx.First(user_yuque_model, user.ID)
		find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
		if find {
			return errors.New(fmt.Sprintf("语雀用户找不到:%d", user.ID))
		}

		//找到用户
		query_result = tx.First(user_model, user_yuque_model.UserId)
		find = errors.Is(query_result.Error, gorm.ErrRecordNotFound)
		if find {
			return errors.New(fmt.Sprintf("用户找不到:%d", user_yuque_model.UserId))
		}

		//更新用户昵称
		user_model.Nickname = user.Name

		result := tx.Save(user_model)

		if result.Error != nil {
			return result.Error
		}

		//更新语雀用户信息
		user_yuque_model.Login = user.Login
		user_yuque_model.Name = user.Name
		user_yuque_model.AvatarUrl = user.AvatarUrl
		user_yuque_model.Description = user.Description
		result = tx.Save(user_model)

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return user_model
}

/**
通过语雀的webhook创建用户
*/
func (s *UserBkService) CreateUserByYuqueWebHook(user *response.UserSerializer) *model.UsereModel {
	db := mysql.GetDefaultDBConnect()

	user_model := new(model.UsereModel)

	err := db.Transaction(func(tx *gorm.DB) error {

		user_yuque_model := new(model.UsereYuQueModel)
		query_result := tx.First(user_yuque_model, user.ID)
		find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
		if !find {
			return errors.New(fmt.Sprintf("语雀用户已同步:%d", user.ID))
		}
		//创建用户
		user_model.Nickname = user.Name
		result := tx.Create(user_model)
		if result.Error != nil {
			return result.Error
		}

		//创建语雀用户
		user_yuque_model.ID = uint(user.ID)
		user_yuque_model.Login = user.Login
		user_yuque_model.Name = user.Name
		user_yuque_model.AvatarUrl = user.AvatarUrl
		user_yuque_model.Description = user.Description
		user_yuque_model.UserId = int64(user_model.ID) //绑定主表
		result = tx.Create(user_yuque_model)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return user_model
}
