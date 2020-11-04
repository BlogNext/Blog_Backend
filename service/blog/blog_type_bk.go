package blog

import (
	"errors"
	"fmt"
	"github.com/FlashFeiFei/yuque/response"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/model"
	"gorm.io/gorm"
)

//博客类型后台服务
type BlogTypeBkService struct {
}

//通过语雀的webhook更新类型
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

//通过语雀webhook创建类型
//book 知识库结构体
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
