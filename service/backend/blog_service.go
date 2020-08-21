package backend

import (
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"time"
)

//博客
type BlogService struct {
}

//添加博客
func (s *BlogService) AddBlog(blog_type_id int64, title, abstract, content string) {
	db := mysql.GetDefaultDBConnect()

	blog_model := new(model.BlogModel)
	blog_model.BlogTypeId = blog_type_id
	blog_model.Title = title
	blog_model.Abstract = abstract
	blog_model.Content = content
	blog_model.CreateTime = time.Now().Unix()
	blog_model.UpdateTime = time.Now().Unix()

	db.Create(blog_model)

	if db.NewRecord(*blog_model) {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("保存失败:%s", db.Error.Error())))
	}
}

//更新博客
func (s *BlogService) UpdateBlog(id, blog_type_id int64, title, abstract, content string) {
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

	db.Save(blog_model)

	if db.Error != nil {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("更新失败error:%s", db.Error.Error())))
	}
}

//列表页
func (s *BlogService) GetList() {

}
