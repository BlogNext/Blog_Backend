package backend

import (
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"strings"
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
func (s *BlogService) GetList() (result []map[string]interface{}) {
	db := mysql.GetDefaultDBConnect()

	//连表表名
	blog_type_table_name := model.BlogTypeModel{}.TableName()
	blog_table_name := model.BlogModel{}.TableName()

	//博客需要的字段
	blog_felid := []string{"id", "blog_type_id", "title", "create_time", "update_time"}

	for index, felid := range blog_felid {
		blog_felid[index] = fmt.Sprintf("%s.%s", blog_table_name, felid)
	}

	//博客类型需要的字段
	blog_type_felid := []string{"title as blog_type_title"}

	for index, felid := range blog_type_felid {
		blog_type_felid[index] = fmt.Sprintf("%s.%s", blog_type_table_name, felid)
	}

	select_felid := append(blog_felid, blog_type_felid...)

	//sql
	rows, _ := db.Table(blog_table_name).
		Joins(fmt.Sprintf("INNER JOIN %s ON %s.blog_type_id = %s.id", blog_type_table_name, blog_table_name, blog_type_table_name)).
		Select(strings.Join(select_felid, ", ")).Rows()

	result = make([]map[string]interface{}, 0)

	for rows.Next() {
		var id int64
		var blog_type_id int64
		var title string
		var create_time int64
		var update_time int64
		var blog_type_title string
		rows.Scan(&id, &blog_type_id, &title, &create_time, &update_time, &blog_type_title)

		item := make(map[string]interface{})
		item["id"] = id
		item["blog_type_id"] = blog_type_id
		item["title"] = title
		item["create_time"] = create_time
		item["update_time"] = update_time
		blog_type_info := make(map[string]interface{})
		blog_type_info["title"] = blog_type_title
		item["blog_type_info"] = blog_type_info
		result = append(result, item)
	}

	return
}
