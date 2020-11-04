package yuque

import (
	"errors"
	"github.com/FlashFeiFei/yuque/response"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/blog"
	user_bk "github.com/blog_backend/service/user"
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
	var user_model *model.UsereModel

	db := mysql.GetDefaultDBConnect()
	user_yuque_model := new(model.UsereYuQueModel)
	query_result := db.First(user_yuque_model, user.ID)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	user_bk_service := new(user_bk.UserBkService)
	if find {
		//创建用户
		user_model = user_bk_service.CreateUserByYuqueWebHook(user)
	} else {
		//更新用户
		user_model = user_bk_service.UpdateUserByYuqueWebHook(user)
	}

	return user_model.ID
}

//同步知识库（博客类型）
func syncBlogType(book *response.BookSerializer) (blog_type_id uint) {
	db := mysql.GetDefaultDBConnect()
	blog_type_model := new(model.BlogTypeModel)
	query_result := db.Where("yuque_id = ?", book.ID).First(blog_type_model)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	blog_type_service := new(blog.BlogTypeBkService)
	if find {
		//找不到博客类型
		blog_type_model = blog_type_service.CreateTypeByYuqueWebHook(book)
	} else {
		//找到博客类型
		blog_type_model = blog_type_service.UpdateTypeByYuqueWebHook(book)
	}

	return blog_type_model.ID
}

//同步博客
func syncBlog(doc *response.DocDetailSerializer, user_id, blog_type_id uint) {
	db := mysql.GetDefaultDBConnect()
	blog_model := new(model.BlogModel)
	query_result := db.First(blog_model, doc.ID)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	blog_service := new(blog.BlogBkService)
	if find {
		//获取博客的封面图和摘要
		blog_service.CreateBlogByYuQueWebHook(doc, user_id, blog_type_id)
	} else {
		//更新文档
		blog_service.UpdateBlogByYuQueWebHook(doc)
	}

}
