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
	db := mysql.GetDefaultDBConnect()
	blog_model := new(model.BlogModel)
	query_result := db.Where("yuque_id = ?", doc.ID).First(blog_model)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("博客未创建yuque_id:%d", doc.ID))
	}

	blog_model.YuqueFormat = doc.Format
	blog_model.YuqueLake = doc.BodyLake
	blog_model.Title = doc.Title
	blog_model.Content = doc.Body

	result := db.Save(blog_model)

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
func (s *BlogBkService) CreateBlogByYuQueWebHook(doc *response.DocDetailSerializer, user_id, blog_type_id uint) {
	db := mysql.GetDefaultDBConnect()
	//查找用户
	user_model := new(model.UserModel)
	query_result := db.First(user_model, user_id)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("找不到用户id:%d", user_id))
	}

	//查找博客分类id
	blog_type_model := new(model.BlogTypeModel)
	query_result = db.First(blog_type_model, blog_type_id)
	find = errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("找不到博客分类blog_type_id:%d", blog_type_id))
	}

	blog_model := new(model.BlogModel)
	query_result = db.Where("yuque_id = ?", doc.ID).First(blog_model)
	find = errors.Is(query_result.Error, gorm.ErrRecordNotFound)
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
	attachment_service := new(attachment.AttachmentRtService)
	attachment_entity_list := attachment_service.DownloadBlogImage(DocIntor.Data.Cover, model.ATTACHMENT_BLOG_Module, model.ATTACHMENT_FILE_TYPE_IMAGE)
	if attachment_entity_list == nil {
		panic("下载封面图失败")
	}

	//创建文档
	//语雀数据
	blog_model.YuqueId = doc.ID
	blog_model.YuqueSlug = doc.Slug
	blog_model.YuqueFormat = doc.Format
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

	//es同步
	//blog_list_entity := ChangeToBlogListEntity(blog_model)
	//blog_es_service := new(BlogEsBkService)
	//blog_es_service.AddDoc(blog_list_entity)

}
