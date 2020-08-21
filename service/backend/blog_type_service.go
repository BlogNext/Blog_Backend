package backend

import (
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"time"
)

//博客类型
type BlogTypeService struct {
}

//获取列表类型接口
func (s *BlogTypeService) List() {

}

//添加
func (s *BlogTypeService) Add(title string) {
	db := mysql.GetDefaultDBConnect()
	blog_type_model := new(model.BlogTypeModel)
	blog_type_model.Title = title
	blog_type_model.CreateTime = time.Now().Unix()
	blog_type_model.UpdateTime = time.Now().Unix()
	db.Create(blog_type_model)

	if db.NewRecord(*blog_type_model) {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("保存失败:%s", db.Error.Error())))
	}
}

func (s *BlogTypeService) Update(id int64, title string) {
	db := mysql.GetDefaultDBConnect()
	blog_type_model := new(model.BlogTypeModel)
	db.Where("id = ?", id).First(blog_type_model)

	if db.NewRecord(*blog_type_model) {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("找不到记录:%d", id)))
	}

	blog_type_model.Title = title

	db.Save(blog_type_model)

	if db.Error != nil {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("更新失败error:%s", db.Error.Error())))
	}
}
