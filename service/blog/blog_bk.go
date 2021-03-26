package blog

import (
	"errors"
	"fmt"
	"github.com/FlashFeiFei/yuque/request/front"
	"github.com/FlashFeiFei/yuque/response"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/attachment"
	"gorm.io/gorm"
	"log"
)

//博客后台服务
type BlogBkService struct {
}

//通过yuquewebhook更新博客
//doc 语雀结构体
func (s *BlogBkService) UpdateBlogByYuQueWebHook(doc *response.DocDetailSerializer) {
	db := mysql.GetNewDB(false)
	blogModel := new(model.BlogModel)
	queryResult := db.Where("yuque_id = ?", doc.ID).First(blogModel)
	find := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("博客未创建yuque_id:%d", doc.ID))
	}

	blogModel.YuqueFormat = doc.Format
	blogModel.YuquePublic = int(doc.Public)
	blogModel.YuqueLake = doc.BodyLake
	blogModel.Title = doc.Title
	blogModel.Content = doc.Body

	result := db.Save(blogModel)

	if result.Error != nil {
		panic(result.Error)
	}

	//es同步
	//blog_list_entity := ChangeToBlogListEntity(blog_model)
	//
	//blog_es_service := new(BlogEsBkService)
	//if blog_model.DocID == "" {
	//	blog_es_service.AddDoc(blog_list_entity)
	//} else {
	//	blog_es_service.UpdateDoc(blog_list_entity)
	//}
}

//通过yuquewebhook创建博客
//doc 语雀结构体
//user_id 用户id
//blog_type_id 博客分类
func (s *BlogBkService) CreateBlogByYuQueWebHook(doc *response.DocDetailSerializer, userId, blogTypeId uint64) {
	db := mysql.GetNewDB(false)
	//查找用户
	userModel := new(model.UserModel)
	queryResult := db.First(userModel, userId)
	find := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("找不到用户id:%d", userId))
	}

	//查找博客分类id
	blogTypeModel := new(model.BlogTypeModel)
	queryResult = db.First(blogTypeModel, blogTypeId)
	find = errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("找不到博客分类blog_type_id:%d", blogTypeId))
	}

	blogModel := new(model.BlogModel)
	queryResult = db.Where("yuque_id = ?", doc.ID).First(blogModel)
	find = errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	if !find {
		log.Println(fmt.Sprintf("博客已存在yuque_id:%d", doc.ID))
		panic(fmt.Sprintf("博客已存在yuque_id:%d", doc.ID))
	}
	log.Println("成功进入创建博客")
	//获取博客的封面图和摘要
	log.Println(doc.Slug, doc.BookId)
	DocIntor := front.GetDocIntorSerializer(doc.Slug, doc.BookId)
	log.Println("下载不了封面图和摘要")
	//下载封面图
	attachmentService := new(attachment.AttachmentRtService)
	attachmentEntityList := attachmentService.DownloadBlogImage(DocIntor.Data.Cover, model.ATTACHMENT_BLOG_Module, model.ATTACHMENT_FILE_TYPE_IMAGE)
	if attachmentEntityList == nil {
		panic("下载封面图失败")
	}

	//创建文档
	//语雀数据
	blogModel.YuqueId = doc.ID
	blogModel.YuqueSlug = doc.Slug
	blogModel.YuqueFormat = doc.Format
	blogModel.YuqueLake = doc.BodyLake
	blogModel.YuquePublic = int(doc.Public)
	blogModel.Title = doc.Title
	blogModel.Content = doc.Body
	blogModel.Abstract = DocIntor.Data.CustomDescription // 摘要
	blogModel.CoverPlanId = attachmentEntityList[0].ID   //封面图
	//系统的数据
	blogModel.UserID = userModel.ID         //用户id
	blogModel.BlogTypeId = blogTypeModel.ID //文章分类id

	result := db.Create(blogModel)
	if result.Error != nil {
		panic(result.Error)
	}

	//es同步
	//blog_list_entity := ChangeToBlogListEntity(blog_model)
	//blog_es_service := new(BlogEsBkService)
	//blog_es_service.AddDoc(blog_list_entity)

}
