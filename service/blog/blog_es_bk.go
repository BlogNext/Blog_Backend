package blog

import (
	"errors"
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/es"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
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

	db := mysql.GetDefaultDBConnect()
	blog_model := new(model.BlogModel)
	query_result := db.First(blog_model, blog_doc.ID)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("博客未创建id:%d", blog_doc.ID))
	}

	result, err := b.BaseEsService.AddDoc(es.BLOG_INDEX, blog_doc)

	if err != nil {
		return nil
	}

	blog_model.DocID = result.Id

	query_result = db.Save(blog_model)

	if query_result.Error != nil {
		panic(query_result.Error)
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

//获取一个文档内容
func (b *BlogEsBkService) GetDocByMysqlId(id uint, blog_doc *blog.BlogEntity) {
	err := b.BaseEsService.GetDocByMysqlId(es.BLOG_INDEX, id, blog_doc)
	if err != nil {
		panic(err)
	}
}
