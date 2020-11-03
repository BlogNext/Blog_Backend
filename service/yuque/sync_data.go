package yuque

import (
	"errors"
	"github.com/FlashFeiFei/yuque/request/front"
	"github.com/FlashFeiFei/yuque/response"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/model"
	"gorm.io/gorm"
)

//webhook数据同步
func SyncData(serializer *response.ResponseDocDetailSerializer) {

	//同步用户
	user_id := syncUserData(serializer.Data.User)

	//同步知识库
	blog_type_id := syncBlogType(serializer.Data.Book)

	//同步博客
	syncBlog(serializer.Data, user_id, blog_type_id)
}

//同步用户
func syncUserData(user *response.UserSerializer) (user_id uint) {
	db := mysql.GetDefaultDBConnect()
	user_model := new(model.UsereModel)
	user_yuque_model := new(model.UsereYuQueModel)
	query_result := db.First(user_yuque_model, user.ID)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		//找不到用户，创建用户

		err := db.Transaction(func(tx *gorm.DB) error {

			//创建用户
			user_model.NickName = user.Name
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
			user_yuque_model.UserId = int64(user_model.ID)
			result = tx.Create(user_yuque_model)
			if result.Error != nil {
				return result.Error
			}

			return nil
		})

		if err != nil {
			panic(err)
		}

	} else {
		//找到用户，更新用户

		//更新用户
		err := db.Transaction(func(tx *gorm.DB) error {

			//更新用户
			query_result = tx.First(user_model, user_yuque_model.UserId)
			find = errors.Is(query_result.Error, gorm.ErrRecordNotFound)
			if find {
				return query_result.Error
			}
			user_model.NickName = user.Name
			result := tx.Save(user_model)
			if result.Error != nil {
				return result.Error
			}

			//更新语雀用户
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

	}

	return user_model.ID
}

//同步知识库（博客类型）
func syncBlogType(book *response.BookSerializer) (blog_type_id uint) {
	db := mysql.GetDefaultDBConnect()
	blog_type_model := new(model.BlogTypeModel)
	query_result := db.Where("yuque_id = ?", book.ID).First(blog_type_model)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		//找不到博客类型
		blog_type_model.YuqueId = book.ID
		blog_type_model.YuqueName = book.Name
		blog_type_model.YuqueType = book.Type

		result := db.Create(blog_type_model)
		if result.Error != nil {
			panic(result.Error)
		}

	} else {
		//找到博客类型
		blog_type_model.YuqueName = book.Name
		blog_type_model.YuqueType = book.Type
		result := db.Save(blog_type_model)

		if result.Error != nil {
			panic(result.Error)
		}
	}

	return blog_type_model.ID
}

//同步博客
func syncBlog(doc *response.DocDetailSerializer, user_id, blog_type_id uint) {
	db := mysql.GetDefaultDBConnect()
	//查找用户
	user_model := new(model.UsereModel)
	query_result := db.First(user_model, user_id)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		panic(query_result.Error)
	}

	//查找博客id
	blog_type_model := new(model.BlogTypeModel)
	query_result = db.First(blog_type_model, blog_type_id)
	find = errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		panic(query_result.Error)
	}

	blog_model := new(model.BlogModel)
	query_result = db.First(blog_model, doc.ID)
	find = errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		//获取博客的封面图和摘要
		DocIntor := front.GetDocIntorSerializer(doc.Slug,doc.BookId)

		//创建文档
		//语雀数据
		blog_model.YuqueId = doc.ID
		blog_model.YuqueSlug = doc.Slug
		blog_model.YuqueIdFormat = doc.Format
		blog_model.YuqueHtml = doc.BodyHtml
		blog_model.YuqueLake = doc.BodyLake
		blog_model.Title = doc.Title
		blog_model.Content = doc.Body
		blog_model.Abstract = DocIntor.Data.CustomDescription

		//系统的数据
		blog_model.UserID = user_model.ID
		blog_model.BlogTypeId = int64(blog_type_model.ID)
	} else {
		//更新文档
	}

}
