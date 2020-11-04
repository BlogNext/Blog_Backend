package blog

import (
	"errors"
	"fmt"
	"github.com/FlashFeiFei/yuque/response"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/model"
	"gorm.io/gorm"
	"log"
	"strings"
)

type BlogTypeBkService struct {
}

/**
通过博客ids获取blogTypeEntity
*/
func (s *BlogTypeBkService) getListByids(ids []uint64) (result map[uint64]*blog.BlogTypeEntity) {
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

/**
获取分类列表
*/
func (s *BlogTypeBkService) GetList(per_page, page int) (result *entity.ListResponseEntity) {
	db := mysql.GetDefaultDBConnect()

	var count int64
	db = db.Model(&model.BlogTypeModel{})
	db.Count(&count)

	rows, _ := db.Select("id, title, create_time, update_time").
		Limit(per_page).Offset((page - 1) * per_page).Rows()

	defer rows.Close()

	blog_type_model_list := make([]map[string]interface{}, 0)

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

	result = new(entity.ListResponseEntity)
	result.SetCount(count)
	result.SetPerPage(per_page)
	result.SetList(blog_type_model_list)

	return
}

/**
通过语雀的webhook更新类型
*/
func (s *BlogTypeBkService) UpdateTypeByYuqueWebHook(book *response.BookSerializer) *model.BlogTypeModel {
	db := mysql.GetDefaultDBConnect()
	blog_type_model := new(model.BlogTypeModel)
	query_result := db.Where("yuque_id = ?", book.ID).First(blog_type_model)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("找不到类型yuque_id:%d", book.ID))
	}

	blog_type_model.YuqueName = book.Name
	blog_type_model.YuqueType = book.Type
	result := db.Save(blog_type_model)

	if result.Error != nil {
		panic(result.Error)
	}

	return blog_type_model
}

/**
通过语雀webhook创建类型
book 知识库结构体
*/
func (s *BlogTypeBkService) CreateTypeByYuqueWebHook(book *response.BookSerializer) *model.BlogTypeModel {
	db := mysql.GetDefaultDBConnect()
	blog_type_model := new(model.BlogTypeModel)
	query_result := db.Where("yuque_id = ?", book.ID).First(blog_type_model)
	find := errors.Is(query_result.Error, gorm.ErrRecordNotFound)
	if !find {
		panic(fmt.Sprintf("已存在类型yuque_id:%d", book.ID))
	}

	blog_type_model.YuqueId = book.ID
	blog_type_model.YuqueName = book.Name
	blog_type_model.YuqueType = book.Type

	result := db.Create(blog_type_model)
	if result.Error != nil {
		panic(result.Error)
	}

	return blog_type_model
}
