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
func (s *BlogTypeService) List() (blog_type_model_list []model.BlogTypeModel) {

	content := mysql.GetDefaultDBConnect()

	db := content.Model(&model.BlogTypeModel{})
	db.Select("id, title, create_time, update_time")
	rows, _ := db.Rows()

	defer rows.Close()

	blog_type_model_list = make([]model.BlogTypeModel, 0)

	for rows.Next() {
		var blog_type_model model.BlogTypeModel
		db.ScanRows(rows, &blog_type_model)
		blog_type_model_list = append(blog_type_model_list, blog_type_model)
	}

	return
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

//修改
func (s *BlogTypeService) Update(id int64, title string) {
	db := mysql.GetDefaultDBConnect()
	blog_type_model := new(model.BlogTypeModel)
	db.Where("id = ?", id).First(blog_type_model)

	if db.NewRecord(*blog_type_model) {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("找不到记录:%d", id)))
	}

	blog_type_model.Title = title
	blog_type_model.UpdateTime = time.Now().Unix()

	db.Save(blog_type_model)

	if db.Error != nil {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("更新失败error:%s", db.Error.Error())))
	}
}
