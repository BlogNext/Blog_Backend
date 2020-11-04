package blog

import (
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/blog"
	"github.com/blog_backend/model"
	"log"
	"strings"
)

//博客前台服务
type BlogRtService struct {
}

//列表页
func (s *BlogRtService) GetList(per_page, page int) (result *entity.ListResponseEntity) {
	db := mysql.GetDefaultDBConnect()

	blog_table_name := model.BlogModel{}.TableName()

	//博客需要的字段
	blog_felid := []string{"id", "blog_type_id", "cover_plan_id", "title", "create_time", "update_time"}

	for index, felid := range blog_felid {
		blog_felid[index] = fmt.Sprintf("%s.%s", blog_table_name, felid)
	}

	var count int64
	//sql
	db = db.Table(blog_table_name)

	db.Count(&count)

	rows, err := db.Select(strings.Join(blog_felid, ", ")).Order("create_time DESC").Limit(per_page).Offset((page - 1) * per_page).Rows()

	if err != nil {
		return nil
	}

	query_result := make([]*blog.BlogEntity, 0)

	cover_plan_ids := make([]uint64, 0)
	blog_type_ids := make([]uint64, 0)

	for rows.Next() {
		var id uint64
		var blog_type_id uint64
		var cover_plan_id uint64
		var title string
		var create_time uint64
		var update_time uint64
		rows.Scan(&id, &blog_type_id, &cover_plan_id, &title, &create_time, &update_time)

		//博客实体
		blog_entity := new(blog.BlogEntity)
		blog_entity.ID = id
		blog_entity.BlogTypeId = blog_type_id
		blog_entity.CoverPlanId = cover_plan_id
		blog_entity.Title = title
		blog_entity.CreateTime = create_time
		blog_entity.UpdateTime = update_time
		log.Println("blog_entity")
		log.Println(blog_entity)

		cover_plan_ids = append(cover_plan_ids, cover_plan_id)
		blog_type_ids = append(blog_type_ids, blog_type_id)

		query_result = append(query_result, blog_entity)
	}

	//填充信息
	PaddingAttachemtInfo(cover_plan_ids, query_result) //填充附件信息
	PaddingBlogTypeInfo(blog_type_ids, query_result)   //博客类型实体

	//构建结果返回
	result = new(entity.ListResponseEntity)
	result.SetCount(count)
	result.SetPerPage(per_page)
	result.SetList(query_result)

	return
}
