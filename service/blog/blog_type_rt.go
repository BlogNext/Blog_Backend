package blog

import (
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/model"
	"strings"
)

//博客类型前台服务
type BlogTypeRtService struct {
}

//通过博客ids获取blogTypeEntity
func (s *BlogTypeRtService) getListByids(ids []uint64) (result map[uint64]*blog.BlogTypeEntity) {
	db := mysql.GetDefaultDBConnect()
	tableName := model.BlogTypeModel{}.TableName()
	selectFelid := []string{"id", "yuque_name", "created_at", "updated_at"}
	rows, err := db.Table(tableName).
		Select(strings.Join(selectFelid, ", ")).Where("id IN (?)", ids).Rows()

	if err != nil {
		return nil
	}

	result = make(map[uint64]*blog.BlogTypeEntity)

	for rows.Next() {
		var id uint64
		var yuqueName string
		var createdAt uint64
		var updatedAt uint64

		rows.Scan(&id, &yuqueName, &createdAt, &updatedAt)

		blogTypeEntity := new(blog.BlogTypeEntity)
		blogTypeEntity.ID = id
		blogTypeEntity.Title = yuqueName
		blogTypeEntity.CreatedAt = createdAt
		blogTypeEntity.UpdatedAt = updatedAt

		result[id] = blogTypeEntity

	}

	return
}

//获取分类列表
func (s *BlogTypeRtService) GetList(perPage, page int) (result *entity.ListResponseEntity) {
	db := mysql.GetDefaultDBConnect()

	var count int64
	db = db.Model(&model.BlogTypeModel{})

	db.Count(&count)

	rows, _ := db.Select("id, yuque_name, yuque_id, created_at, updated_at").
		Limit(perPage).Offset((page - 1) * perPage).Rows()

	defer rows.Close()

	blogTypeModelList := make([]map[string]interface{}, 0)

	for rows.Next() {

		var id int64
		var yuqueName string
		var yuqueId int64
		var createdAt int64
		var updatedAt int64
		rows.Scan(&id, &yuqueName, &yuqueId, &createdAt, &updatedAt)

		item := make(map[string]interface{})

		item["id"] = id
		item["yuque_name"] = yuqueName
		item["yuque_id"] = yuqueId
		item["create_time"] = createdAt
		item["update_time"] = updatedAt

		blogTypeModelList = append(blogTypeModelList, item)
	}

	result = new(entity.ListResponseEntity)
	result.SetCount(count)
	result.SetPerPage(perPage)
	result.SetList(blogTypeModelList)

	return
}
