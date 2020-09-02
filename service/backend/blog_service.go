package backend

import (
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity/attachment"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/common"
	es_blog "github.com/blog_backend/service/common/es/blog"
	"github.com/thoas/go-funk"
	"log"
	"strings"
	"time"
)

//博客
type BlogService struct {
}

//模型转化成BlogEntity实体
func (s *BlogService) changeToBlogEntity(blog_model *model.BlogModel) *blog.BlogEntity {

	blog_entity := new(blog.BlogEntity)
	blog_entity.ID = uint64(blog_model.ID)
	blog_entity.CreateTime = uint64(blog_model.CreateTime)
	blog_entity.UpdateTime = uint64(blog_model.UpdateTime)
	blog_entity.BlogTypeId = uint64(blog_model.BlogTypeId)
	blog_entity.CoverPlanId = uint64(blog_model.CoverPlanId)
	blog_entity.Title = blog_model.Title
	blog_entity.Abstract = blog_model.Abstract
	blog_entity.Content = blog_model.Content
	blog_entity.DocID = blog_model.DocID

	log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", blog_entity, blog_entity, blog_entity))

	blog_entity_list := []*blog.BlogEntity{blog_entity}
	//填充别的实体信息

	log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", blog_entity_list, blog_entity_list, blog_entity_list))

	s.paddingAttachemtInfo([]uint64{blog_entity.CoverPlanId}, blog_entity_list)

	log.Println("附件填充完毕", fmt.Sprintf("v = %v,t = %T, p = %p", blog_entity_list, blog_entity_list, blog_entity_list))

	s.paddingBlogTypeInfo([]uint64{blog_entity.BlogTypeId}, blog_entity_list)

	log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", blog_entity_list, blog_entity_list, blog_entity_list))

	return blog_entity_list[0]
}

//导入数据到es中
func (s *BlogService) ImportDataToEs() {

	var blog_list []model.BlogModel

	db := mysql.GetDefaultDBConnect()
	db.Find(&blog_list)

	if blog_list == nil {
		return
	}

	log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", blog_list, blog_list, blog_list))

	for _, blog_model := range blog_list {
		//es中添加文件
		blog_doc := s.changeToBlogEntity(&blog_model)

		log.Println("导入的es文档是：", fmt.Sprintf("v = %v,t = %T, p = %p", blog_doc, blog_doc, blog_doc))

		es_blog_service, err := es_blog.NewBlogEsService("", "", "")

		if err != nil {
			log.Println("连接es失败结束同步")
			break
		}
		log.Println("连接:es成功")

		doc := es_blog_service.AddDoc(blog_doc)

		blog_model.DocID = doc.Id

		db_error := db.Save(blog_model)

		if db_error.Error != nil {
			panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("更新失败error:%s", db_error.Error.Error())))
		}
	}

}

//添加博客
func (s *BlogService) AddBlog(blog_type_id, cover_plan_id int64, title, abstract, content string) {

	//数据入库
	db := mysql.GetDefaultDBConnect()

	blog_model := new(model.BlogModel)
	blog_model.BlogTypeId = blog_type_id
	blog_model.Title = title
	blog_model.Abstract = abstract
	blog_model.Content = content
	blog_model.CoverPlanId = cover_plan_id
	blog_model.CreateTime = time.Now().Unix()
	blog_model.UpdateTime = time.Now().Unix()

	db.Create(blog_model)

	if db.NewRecord(*blog_model) {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("保存失败:%s", db.Error.Error())))
	}

	//创建es文档
	blog_doc := s.changeToBlogEntity(blog_model) //文档转化

	es_blog_service, err := es_blog.NewBlogEsService("", "", "")

	if err != nil {
		log.Println("连接es失败结束同步，es入库失败")
		return
	}

	doc := es_blog_service.AddDoc(blog_doc)

	blog_model.DocID = doc.Index   //文档保存

	db_error := db.Save(blog_model)

	if db_error.Error != nil {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("更新失败error:%s", db_error.Error.Error())))
	}

}

//更新博客
func (s *BlogService) UpdateBlog(id, blog_type_id, cover_plan_id int64, title, abstract, content string) {
	db := mysql.GetDefaultDBConnect()
	blog_model := new(model.BlogModel)
	db.Where("id = ?", id).First(blog_model)

	if db.NewRecord(*blog_model) {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("找不到记录:%d", id)))
	}

	blog_model.BlogTypeId = blog_type_id
	blog_model.Title = title
	blog_model.Abstract = abstract
	blog_model.Content = content
	blog_model.UpdateTime = time.Now().Unix()
	blog_model.CoverPlanId = cover_plan_id

	db.Save(blog_model)

	if db.Error != nil {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("更新失败error:%s", db.Error.Error())))
	}

	//更新es文档
	blog_doc := s.changeToBlogEntity(blog_model) //文档转化

	es_blog_service, err := es_blog.NewBlogEsService("", "", "")

	if err != nil {
		log.Println("连接es失败结束同步，es入库失败")
		return
	}

	_ = es_blog_service.UpdateDoc(blog_doc)



}

//列表页
func (s *BlogService) GetList() (result []*blog.BlogEntity) {
	db := mysql.GetDefaultDBConnect()

	blog_table_name := model.BlogModel{}.TableName()

	//博客需要的字段
	blog_felid := []string{"id", "blog_type_id", "cover_plan_id", "title", "create_time", "update_time"}

	for index, felid := range blog_felid {
		blog_felid[index] = fmt.Sprintf("%s.%s", blog_table_name, felid)
	}

	//sql
	rows, _ := db.Table(blog_table_name).Select(strings.Join(blog_felid, ", ")).Rows()

	result = make([]*blog.BlogEntity, 0)

	cover_plan_ids := make([]uint64, 0)
	blog_type_ids := make([]uint64, 0)

	for rows.Next() {
		var id uint64
		var blog_type_id uint64
		var cover_plan_id uint64
		var title string
		var create_time uint64
		var update_time uint64
		rows.Scan(&id, &blog_type_id, &cover_plan_id, &title, &create_time, &update_time)

		//博客实体
		blog_entity := new(blog.BlogEntity)
		blog_entity.ID = id
		blog_entity.BlogTypeId = blog_type_id
		blog_entity.CoverPlanId = cover_plan_id
		blog_entity.Title = title
		blog_entity.CreateTime = create_time
		blog_entity.UpdateTime = update_time

		cover_plan_ids = append(cover_plan_ids, cover_plan_id)
		blog_type_ids = append(blog_type_ids, blog_type_id)

		result = append(result, blog_entity)
	}

	//填充信息
	s.paddingAttachemtInfo(cover_plan_ids, result) //填充附件信息
	s.paddingBlogTypeInfo(blog_type_ids, result)   //博客类型实体

	return
}

//填充附件信息
func (s *BlogService) paddingAttachemtInfo(cover_plan_ids []uint64, result []*blog.BlogEntity) {

	//获取图片的ids,填充图片信息
	attachment_list := common.GetAttachmentImages(cover_plan_ids)

	if attachment_list != nil {
		//转化成map
		log.Println(attachment_list)
		log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", attachment_list, attachment_list, attachment_list))

		attachment_list_map := funk.ToMap(attachment_list, "ID").(map[uint64]*attachment.AttachmentEntity)
		log.Println(attachment_list_map)
		log.Println("哗啦啦")
		//填充图片信息
		for _, item := range result {
			log.Println(attachment_list_map[item.CoverPlanId])

			if attachment_item, ok := attachment_list_map[item.CoverPlanId]; ok {
				item.AttachmentInfo = attachment_item
			} else {
				item.AttachmentInfo = nil
			}

		}
	}
}

//填充博客类型信息
func (s *BlogService) paddingBlogTypeInfo(blog_type_ids []uint64, result []*blog.BlogEntity) {

	blog_type_service := new(BlogTypeService)
	blog_type_list := blog_type_service.getListByids(blog_type_ids)

	log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", blog_type_list, blog_type_list, blog_type_list))

	if blog_type_list != nil {
		//填充博客类型

		for _, item := range result {
			if blog_type_entity, ok := blog_type_list[item.BlogTypeId]; ok {
				item.BlogTypeObject = blog_type_entity
			} else {
				item.BlogTypeObject = nil
			}
		}
	}
}
