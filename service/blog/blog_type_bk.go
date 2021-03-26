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
	db := mysql.GetNewDB(false)
	blogTypeModel := new(model.BlogTypeModel)
	queryResult := db.Where("yuque_id = ?", book.ID).First(blogTypeModel)
	find := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	if find {
		panic(fmt.Sprintf("找不到类型yuque_id:%d", book.ID))
	}

	blogTypeModel.YuqueName = book.Name
	blogTypeModel.YuqueType = book.Type
	result := db.Save(blogTypeModel)

	if result.Error != nil {
		panic(result.Error)
	}

	return blogTypeModel
}

//通过语雀webhook创建类型
//book 知识库结构体
func (s *BlogTypeBkService) CreateTypeByYuqueWebHook(book *response.BookSerializer) *model.BlogTypeModel {
	db := mysql.GetNewDB(false)
	blogTypeModel := new(model.BlogTypeModel)
	queryResult := db.Where("yuque_id = ?", book.ID).First(blogTypeModel)
	find := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	if !find {
		panic(fmt.Sprintf("已存在类型yuque_id:%d", book.ID))
	}

	blogTypeModel.YuqueId = book.ID
	blogTypeModel.YuqueName = book.Name
	blogTypeModel.YuqueType = book.Type

	result := db.Create(blogTypeModel)
	if result.Error != nil {
		panic(result.Error)
	}

	return blogTypeModel
}
