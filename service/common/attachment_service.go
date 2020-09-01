package common

import (
	"fmt"
	"github.com/blog_backend/common-lib/config"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity/attachment"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	//上传的更目录
	UPLOAD_ROOT_PATH = "upload"
)

//获取附件
func GetAttachmentImages(ids []uint64) (attachment_entity_list []*attachment.AttachmentEntity) {
	attachment_list := getAttachmentByIds(ids)
	if attachment_list == nil {
		return
	}

	log.Println(fmt.Sprintf("附件长度=%d", len(attachment_list)))
	attachment_entity_list = make([]*attachment.AttachmentEntity, len(attachment_list))

	server_config, _ := config.GetConfig("server")
	server_info := server_config.GetStringMap("servier")
	domain := server_info["domain"].(string)

	for index, attachment_model := range attachment_list {
		attachment_entity := new(attachment.AttachmentEntity)
		attachment_entity.ID = uint64(attachment_model.ID)
		attachment_entity.CreateTime = uint64(attachment_model.CreateTime)
		attachment_entity.UpdateTime = uint64(attachment_model.UpdateTime)
		attachment_entity.Module = attachment_model.Module
		attachment_entity.Path = attachment_model.Path
		attachment_entity.Url = attachment_model.Path
		attachment_entity.FullUrl = strings.Join([]string{domain, attachment_model.Path}, "/")
		log.Println(attachment_entity)

		attachment_entity_list[index] = attachment_entity
	}

	log.Println(fmt.Sprintf("v = %v,t = %T, p = %p", attachment_entity_list[0], attachment_entity_list[0], attachment_entity_list[0]))

	return
}

//获取附件
func getAttachmentByIds(ids []uint64) (attachment_list []model.AttachmentModel) {

	if ids == nil {
		return
	}

	db := mysql.GetDefaultDBConnect()
	db.Where("id IN (?)", ids).Find(&attachment_list)

	return
}

//附件服务
type AttachmentService struct {
}

//保存到数据库
func (s *AttachmentService) saveToDB(dst string, module int64) (attachment_model *model.AttachmentModel) {

	db := mysql.GetDefaultDBConnect()
	attachment_model = new(model.AttachmentModel)
	attachment_model.Path = dst
	attachment_model.CreateTime = time.Now().Unix()
	attachment_model.UpdateTime = time.Now().Unix()
	attachment_model.Module = module

	db.Create(attachment_model)

	if db.NewRecord(*attachment_model) {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("新增失败:%s", db.Error.Error())))
	}

	return attachment_model
}

//上传博客的文件
func (s *AttachmentService) UploadBlog(Ctx *gin.Context) (full_attachment_extend []*attachment.AttachmentEntity) {
	multipart_form, _ := Ctx.MultipartForm()
	files := multipart_form.File["upload_blog_images"]

	rand.Seed(time.Now().UnixNano())

	dir := strings.Join([]string{UPLOAD_ROOT_PATH, "blog"}, "/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	var attachment_ids []uint64

	//保存成功的文件
	attachment_ids = make([]uint64, len(files))

	for index, file := range files {

		file_rename := fmt.Sprintf("%d.%d.%d", time.Now().UnixNano(), rand.Int63(), index)

		dst := strings.Join([]string{dir, file_rename}, "/")

		err := Ctx.SaveUploadedFile(file, dst)
		if err != nil {
			panic(err)
		}

		//保存到数据库
		attachment_model := s.saveToDB(dst, model.ATTACHMENT_BLOG_Module)

		attachment_ids = append(attachment_ids, uint64(attachment_model.ID))
	}

	//获取文件列表
	full_attachment_extend = GetAttachmentImages(attachment_ids)

	return
}
