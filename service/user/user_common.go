package user

import (
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity/user"
	"github.com/blog_backend/model"
)

//通过用户ids获取UserModel
func GetUserByUserIds(ids []uint64) (userModelList []*model.UserModel) {
	if ids == nil {
		panic("ids不能为空")
	}

	db := mysql.GetNewDB(false)
	db = db.Table(model.UserModel{}.TableName())
	db.Where("id IN ?", ids).Find(&userModelList)

	return userModelList
}

//通过用户ids获取UserEntity
func GetUserEntityByUserIds(ids []uint64) (userEntityList map[uint64]*user.UserEntity) {
	userModelList := GetUserByUserIds(ids)
	if userModelList == nil || len(userModelList) <= 0 {
		return nil
	}

	userEntityList = make(map[uint64]*user.UserEntity, len(userModelList))

	for _, userModel := range userModelList {
		userEntity := new(user.UserEntity)
		userEntity.ID = userModel.ID
		userEntity.Nickname = userModel.Nickname

		userEntityList[userModel.ID] = userEntity
	}

	return userEntityList
}
