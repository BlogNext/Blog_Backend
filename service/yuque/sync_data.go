package yuque

import (
	"errors"
	"github.com/FlashFeiFei/yuque/request/front"
	"github.com/FlashFeiFei/yuque/response"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/attachment"
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
		DocIntor := front.GetDocIntorSerializer(doc.Slug, doc.BookId)
		//下载封面图
		attachment_service := new(attachment.AttachmentService)
		attachment_entity_list := attachment_service.DownloadBlogImage(DocIntor.Data.Cover, model.ATTACHMENT_BLOG_Module, model.ATTACHMENT_FILE_TYPE_IMAGE)
		if attachment_entity_list == nil {
			panic("下载封面图失败")
		}
		//创建文档
		//语雀数据
		blog_model.YuqueId = doc.ID
		blog_model.YuqueSlug = doc.Slug
		blog_model.YuqueIdFormat = doc.Format
		blog_model.YuqueHtml = doc.BodyHtml
		blog_model.YuqueLake = doc.BodyLake
		blog_model.Title = doc.Title
		blog_model.Content = doc.Body
		blog_model.Abstract = DocIntor.Data.CustomDescription        // 摘要
		blog_model.CoverPlanId = int64(attachment_entity_list[0].ID) //封面图
		//系统的数据
		blog_model.UserID = user_model.ID                 //用户id
		blog_model.BlogTypeId = int64(blog_type_model.ID) //文章分类id

		result := db.Create(blog_model)
		if result.Error != nil {
			panic(result.Error)
		}

	} else {
		//更新文档
		blog_model.YuqueIdFormat = doc.Format
		blog_model.YuqueHtml = doc.BodyHtml
		blog_model.YuqueLake = doc.BodyLake
		blog_model.Title = doc.Title
		blog_model.Content = doc.Body

		result := db.Save(blog_model)

		if result.Error != nil {
			panic(result.Error)
		}
	}

}
