package blog

import (
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/model"
	"log"
	"strings"
)

//博客类型前台服务
type BlogTypeRtService struct {
}

//通过博客ids获取blogTypeEntity
func (s *BlogTypeRtService) getListByids(ids []uint64) (result map[uint64]*blog.BlogTypeEntity) {
	db := mysql.GetDefaultDBConnect()
	table_name := model.BlogTypeModel{}.TableName()
	select_felid := []string{"id", "title", "created_at", "updated_at"}
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
		var created_at uint64
		var updated_at uint64

		rows.Scan(&id, &title, &created_at, &updated_at)

		blog_type_entity := new(blog.BlogTypeEntity)
		blog_type_entity.ID = id
		blog_type_entity.Title = title
		blog_type_entity.CreatedAt = created_at
		blog_type_entity.UpdatedAt = updated_at

		result[id] = blog_type_entity

	}

	return
}

//获取分类列表
func (s *BlogTypeRtService) GetList(per_page, page int) (result *entity.ListResponseEntity) {
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
