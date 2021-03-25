package user

import (
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity/user"
	"github.com/blog_backend/model"
)

//通过用户ids获取UserModel
func GetUserByUserIds(ids []uint64) (user_model_list []*model.UserModel) {
	if ids == nil {
		panic("ids不能为空")
	}

	db := mysql.GetDefaultDBConnect()
	db = db.Table(model.UserModel{}.TableName())
	db.Where("id IN ?", ids).Find(&user_model_list)

	return user_model_list
}

//通过用户ids获取UserEntity
func GetUserEntityByUserIds(ids []uint64) (user_entity_list map[uint64]*user.UserEntity) {
	user_model_list := GetUserByUserIds(ids)
	if user_model_list == nil || len(user_model_list) <= 0 {
		return nil
	}

	user_entity_list = make(map[uint64]*user.UserEntity, len(user_model_list))

	for _, user_model := range user_model_list {
		user_entity := new(user.UserEntity)
		user_entity.ID = uint64(user_model.ID)
		user_entity.Nickname = user_model.Nickname

		user_entity_list[user_model.ID] = user_entity
	}

	return user_entity_list
}
