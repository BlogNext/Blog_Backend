package yuque

import (
	"encoding/json"
	"errors"
	"github.com/FlashFeiFei/yuque/request"
	"github.com/FlashFeiFei/yuque/response"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/blog"
	user_bk "github.com/blog_backend/service/user"
	"gorm.io/gorm"
	"log"
)

//webhook数据同步
func SyncData(serializer *response.ResponseDocDetailSerializer, token string) {

	log.Println("语雀Token: ", token)
	//同步用户
	userRequest := request.UserRequest{AuthToken: request.AuthToken{Token: token}} //获取文章人创建人信息
	userInfoRequest := userRequest.NewUserRequestById(serializer.Data.UserId)
	userResponse := new(response.ResponseUserSerializer)
	userInfoRequest.Request(userResponse)
	log.Println("用户信息")
	log.Println(json.Marshal(userInfoRequest))
	userId := syncUserData(userResponse.Data)

	//同步知识库
	blogTypeId := syncBlogType(serializer.Data.Book)

	//同步博客
	syncBlog(serializer.Data, userId, blogTypeId)
}

//同步用户
func syncUserData(user *response.UserSerializer) (userId uint64) {
	var userModel *model.UserModel

	db := mysql.GetNewDB(false)
	userYuqueModel := new(model.UserYuQueModel)
	queryResult := db.First(userYuqueModel, user.ID)
	find := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	userBkService := new(user_bk.UserBkService)
	if find {
		//创建用户
		userModel = userBkService.CreateUserByYuqueWebHook(user)
	} else {
		//更新用户
		userModel = userBkService.UpdateUserByYuqueWebHook(user)
	}

	return userModel.ID
}

//同步知识库（博客类型）
func syncBlogType(book *response.BookSerializer) (blogTypeId uint64) {
	db := mysql.GetNewDB(false)
	blogTypeModel := new(model.BlogTypeModel)
	queryResult := db.Where("yuque_id = ?", book.ID).First(blogTypeModel)
	find := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	blogTypeService := new(blog.BlogTypeBkService)
	if find {
		//找不到博客类型
		blogTypeModel = blogTypeService.CreateTypeByYuqueWebHook(book)
	} else {
		//找到博客类型
		blogTypeModel = blogTypeService.UpdateTypeByYuqueWebHook(book)
	}

	return blogTypeModel.ID
}

//同步博客
func syncBlog(doc *response.DocDetailSerializer, userId, blogTypeId uint64) {
	db := mysql.GetNewDB(false)
	blogModel := new(model.BlogModel)
	queryResult := db.Where("yuque_id = ?", doc.ID).First(blogModel)
	find := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	blogService := new(blog.BlogBkService)
	if find {
		//获取博客的封面图和摘要
		blogService.CreateBlogByYuQueWebHook(doc, userId, blogTypeId)
	} else {
		//更新文档
		blogService.UpdateBlogByYuQueWebHook(doc)
	}

}
