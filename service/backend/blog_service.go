package backend

import (
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity"
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
		blog_doc := s.ChangeToBlogEntity(&blog_model)

		log.Println("导入的es文档是：", fmt.Sprintf("v = %v,t = %T, p = %p", blog_doc, blog_doc, blog_doc))

		es_blog_service := new(es_blog.BlogEsService)

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

	sql_exec_result := db.Create(blog_model)

	if sql_exec_result.Error != nil {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("新增失败:%s", sql_exec_result.Error)))
	}

	//创建es文档
	blog_doc := s.ChangeToBlogEntity(blog_model) //文档转化

	es_blog_service := new(es_blog.BlogEsService)

	doc := es_blog_service.AddDoc(blog_doc)

	blog_model.DocID = doc.Id //文档保存

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

	if blog_model.ID <= 0 {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("找不到记录:%d", id)))
	}

	blog_model.BlogTypeId = blog_type_id
	blog_model.Title = title
	blog_model.Abstract = abstract
	blog_model.Content = content
	blog_model.UpdateTime = time.Now().Unix()
	if cover_plan_id != 0 {
		blog_model.CoverPlanId = cover_plan_id
	}

	log.Println(fmt.Sprintf("更新的博客数据: v= %v, t= %T, p=%p", blog_model, blog_model, blog_model))
	result := db.Save(blog_model)

	if result.Error != nil {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("更新失败error:%s", result.Error.Error())))
	}

	//更新es文档
	blog_doc := s.ChangeToBlogEntity(blog_model) //文档转化

	log.Println("转化为es的文档为", fmt.Sprintf("v=%v, t=%T ,p=%p", blog_doc, blog_doc, blog_doc))

	es_blog_service := new(es_blog.BlogEsService)

	_ = es_blog_service.UpdateDoc(blog_doc)

}

