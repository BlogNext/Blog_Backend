package blog

import (
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/es"
	"github.com/olivere/elastic/v7"
	"log"
)

type BlogEsBkService struct {
	es.BaseEsService
}

//创建博客索引
func (b *BlogEsBkService) CreateIndex() {
	commend := es.CreateBlogIndex()
	//设置命令
	b.SetExecCommend(commend)
	//运行命令
	_, err := es.RunCommend(b)
	if err != nil {
		log.Println("无法创建索引")
		panic(err)
	}
}

//导入数据到es中
func (s *BlogEsBkService) ImportDataToEs() {

	var blog_list []model.BlogModel

	db := mysql.GetDefaultDBConnect()
	db.Find(&blog_list)

	if blog_list == nil {
		return
	}

	log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", blog_list, blog_list, blog_list))

	for _, blog_model := range blog_list {
		//es中添加文件
		blog_doc := ChangeToBlogEntity(&blog_model)

		log.Println("导入的es文档是：", fmt.Sprintf("v = %v,t = %T, p = %p", blog_doc, blog_doc, blog_doc))

		es_blog_service := new(BlogEsBkService)

		log.Println("连接:es成功")

		doc := es_blog_service.AddDoc(blog_doc)

		blog_model.DocID = doc.Id

		db_error := db.Save(blog_model)

		if db_error.Error != nil {
			panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("更新失败error:%s", db_error.Error.Error())))
		}
	}

}

//删除blog文档
func (b *BlogEsBkService) DeleteDoc(blog_doc *blog.BlogEntity) *elastic.DeleteResponse {
	result, err := b.BaseEsService.DeleteDoc(es.BLOG_INDEX, blog_doc.DocID)
	if err != nil {
		return nil
	}
	return result
}

//添加一个doc,返回文档在es中的唯一标识
func (b *BlogEsBkService) AddDoc(blog_doc *blog.BlogEntity) *elastic.IndexResponse {
	result, err := b.BaseEsService.AddDoc(es.BLOG_INDEX, blog_doc)
	if err != nil {
		return nil
	}
	return result
}

//更新一个文档的内容
func (b *BlogEsBkService) UpdateDoc(blog_doc *blog.BlogEntity) *elastic.UpdateResponse {
	result, err := b.BaseEsService.UpdateDoc(es.BLOG_INDEX, blog_doc.DocID, blog_doc)
	if err != nil {
		return nil
	}
	return result
}
