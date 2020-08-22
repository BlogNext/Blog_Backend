package common

import (
	"fmt"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"github.com/gin-gonic/gin"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	//上传的更目录
	UPLOAD_ROOT_PATH = "upload"
)

type AttachmentService struct {
}

//保存到数据库
func (s *AttachmentService) saveToDB(dst string, module int64) {

	db := mysql.GetDefaultDBConnect()
	attachment_model := new(model.AttachmentModel)
	attachment_model.Path = dst
	attachment_model.CreateTime = time.Now().Unix()
	attachment_model.UpdateTime = time.Now().Unix()
	attachment_model.Module = module

	db.Create(attachment_model)

	if db.NewRecord(*attachment_model) {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("新增失败:%s", db.Error.Error())))
	}
}

//上传博客的文件
func (s *AttachmentService) UploadBlog(Ctx *gin.Context) {
	multipart_form, _ := Ctx.MultipartForm()
	files := multipart_form.File["upload_blog_images"]

	rand.Seed(time.Now().UnixNano())

	dir := strings.Join([]string{UPLOAD_ROOT_PATH, "blog"}, "/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	for index, file := range files {

		file_rename := fmt.Sprintf("%d.%d.%d", time.Now().UnixNano(), rand.Int63(), index)

		dst := strings.Join([]string{dir, file_rename}, "/")

		err := Ctx.SaveUploadedFile(file, dst)
		if err != nil {
			panic(err)
		}

		//保存到数据库
		s.saveToDB(dst, model.ATTACHMENT_BLOG_Module)
	}

}
