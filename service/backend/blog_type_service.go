package backend

import (
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"log"
	"strings"
	"time"
)

//博客类型
type BlogTypeService struct {
}

//获取列表类型接口
func (s *BlogTypeService) List() (blog_type_model_list []map[string]interface{}) {

	content := mysql.GetDefaultDBConnect()

	db := content.Model(&model.BlogTypeModel{})
	db.Select("id, title, create_time, update_time")
	rows, _ := db.Rows()

	defer rows.Close()

	blog_type_model_list = make([]map[string]interface{}, 0)

	for rows.Next() {

		var id int64
		var title string
		var create_time int64
		var update_time int64
		rows.Scan(&id, &title, &create_time, &update_time)

		item := make(map[string]interface{})

		item["id"] = id
		item["title"] = title
		item["create_time"] = create_time
		item["update_time"] = update_time

		blog_type_model_list = append(blog_type_model_list, item)
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

//通过ids获取类型信息
//返回以id为key，值为BlogTypeEntity
func (s *BlogTypeService) getListByids(ids []uint64) (result map[uint64]*blog.BlogTypeEntity) {

	log.Println("到这里")

	log.Println("到这里", ids)
	db := mysql.GetDefaultDBConnect()
	table_name := model.BlogTypeModel{}.TableName()
	select_felid := []string{"id", "title", "create_time", "update_time"}
	rows, err := db.Table(table_name).
		Select(strings.Join(select_felid, ", ")).Where("id IN (?)", ids).Rows()

	if err != nil {
		return nil
	}

	log.Println("到这里", rows)
	result = make(map[uint64]*blog.BlogTypeEntity)

	for rows.Next() {
		var id uint64
		var title string
		var create_time uint64
		var update_time uint64

		rows.Scan(&id, &title, &create_time, &update_time)

		blog_type_entity := new(blog.BlogTypeEntity)
		blog_type_entity.ID = id
		blog_type_entity.Title = title
		blog_type_entity.CreateTime = create_time
		blog_type_entity.UpdateTime = update_time

		result[id] = blog_type_entity

	}

	return
}
