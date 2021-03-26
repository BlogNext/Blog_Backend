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

func (s *UserBkService) UpdateUserByYuqueWebHook(user *response.UserSerializer) *model.UserModel {
	db := mysql.GetNewDB(false)
	userModel := new(model.UserModel)

	//更新用户
	err := db.Transaction(func(tx *gorm.DB) error {

		//找到语雀用户
		userYuqueModel := new(model.UserYuQueModel)
		query_result := tx.First(userYuqueModel, user.ID)
		find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
		if find {
			return errors.New(fmt.Sprintf("语雀用户找不到:%d", user.ID))
		}

		//找到用户
		query_result = tx.First(userModel, userYuqueModel.UserId)
		find = errors.Is(query_result.Error, gorm.ErrRecordNotFound)
		if find {
			return errors.New(fmt.Sprintf("用户找不到:%d", userYuqueModel.UserId))
		}

		//更新用户昵称
		userModel.Nickname = user.Name

		result := tx.Save(userModel)

		if result.Error != nil {
			return result.Error
		}

		//更新语雀用户信息
		userYuqueModel.Login = user.Login
		userYuqueModel.Name = user.Name
		userYuqueModel.AvatarUrl = user.AvatarUrl
		userYuqueModel.Description = user.Description
		result = tx.Save(userModel)

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return userModel
}

/**
通过语雀的webhook创建用户
*/
func (s *UserBkService) CreateUserByYuqueWebHook(user *response.UserSerializer) *model.UserModel {
	db := mysql.GetNewDB(false)

	userModel := new(model.UserModel)

	err := db.Transaction(func(tx *gorm.DB) error {

		userYuqueModel := new(model.UserYuQueModel)
		queryResult := tx.First(userYuqueModel, user.ID)
		find := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
		if !find {
			return errors.New(fmt.Sprintf("语雀用户已同步:%d", user.ID))
		}
		//创建用户
		userModel.Nickname = user.Name
		result := tx.Create(userModel)
		if result.Error != nil {
			return result.Error
		}

		//创建语雀用户
		userYuqueModel.ID = uint64(user.ID)
		userYuqueModel.Login = user.Login
		userYuqueModel.Name = user.Name
		userYuqueModel.AvatarUrl = user.AvatarUrl
		userYuqueModel.Description = user.Description
		userYuqueModel.UserId = userModel.ID //绑定主表
		result = tx.Create(userYuqueModel)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return userModel
}
