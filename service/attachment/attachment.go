package attachment

import (
	"fmt"
	"github.com/blog_backend/common-lib/config"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity/attachment"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

const (
	//上传的更目录
	UPLOAD_ROOT_PATH = "upload"
)

//获取附件
//返回map
func GetAttachmentImagesMap(ids []uint64) (attachment_entity_list map[uint]*attachment.AttachmentEntity) {
	attachment_list := getAttachmentByIds(ids)
	if attachment_list == nil {
		return
	}

	if len(attachment_list) <= 0 {
		return
	}

	attachment_entity_list = make(map[uint]*attachment.AttachmentEntity, len(attachment_list))

	server_config, _ := config.GetConfig("server")
	server_info := server_config.GetStringMap("servier")
	fileDomain := server_info["fileDomain"].(string)

	for _, attachment_model := range attachment_list {
		attachment_entity := new(attachment.AttachmentEntity)
		attachment_entity.ID = uint64(attachment_model.ID)
		attachment_entity.CreatedAt = uint64(attachment_model.CreatedAt)
		attachment_entity.UpdatedAt = uint64(attachment_model.UpdatedAt)
		attachment_entity.Module = attachment_model.Module
		attachment_entity.Path = attachment_model.Path
		attachment_entity.Url = attachment_model.Path
		attachment_entity.FullUrl = strings.Join([]string{fileDomain, attachment_model.Path}, "/")
		attachment_entity.FileType = attachment_model.FileType

		attachment_entity_list[uint(attachment_entity.ID)] = attachment_entity
	}

	return
}

//获取附件
//返回切片
func GetAttachmentImages(ids []uint64) (attachment_entity_list []*attachment.AttachmentEntity) {
	attachment_list := getAttachmentByIds(ids)
	if attachment_list == nil {
		return
	}

	if len(attachment_list) <= 0 {
		return
	}

	attachment_entity_list = make([]*attachment.AttachmentEntity, len(attachment_list))

	server_config, _ := config.GetConfig("server")
	server_info := server_config.GetStringMap("servier")
	domain := server_info["domain"].(string)

	for index, attachment_model := range attachment_list {
		attachment_entity := new(attachment.AttachmentEntity)
		attachment_entity.ID = uint64(attachment_model.ID)
		attachment_entity.CreatedAt = uint64(attachment_model.CreatedAt)
		attachment_entity.UpdatedAt = uint64(attachment_model.UpdatedAt)
		attachment_entity.Module = attachment_model.Module
		attachment_entity.Path = attachment_model.Path
		attachment_entity.Url = attachment_model.Path
		attachment_entity.FullUrl = strings.Join([]string{domain, attachment_model.Path}, "/")
		attachment_entity.FileType = attachment_model.FileType

		attachment_entity_list[index] = attachment_entity
	}

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
type AttachmentBaseService struct {
}

//保存到数据库
func (s *AttachmentBaseService) saveToDB(dst string, module, file_type int64) (attachment_model *model.AttachmentModel) {

	db := mysql.GetDefaultDBConnect()
	attachment_model = new(model.AttachmentModel)
	attachment_model.Path = dst
	attachment_model.CreatedAt = time.Now().Unix()
	attachment_model.UpdatedAt = time.Now().Unix()
	attachment_model.Module = module
	attachment_model.FileType = file_type

	sql_exec_result := db.Create(attachment_model)

	if sql_exec_result.Error != nil {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("新增失败:%s", sql_exec_result.Error)))
	}

	return attachment_model
}

//重命名文件名
//file_name_list 一批文件名
func (s *AttachmentBaseService) renameFileName(file_name_list []string) (new_file_name_list []string) {

	rand.Seed(time.Now().UnixNano())

	if file_name_list == nil || len(file_name_list) == 0 {
		return nil
	}

	new_file_name_list = make([]string, len(file_name_list))
	for index, file_name := range file_name_list {
		new_file_name_list[index] = fmt.Sprintf("%d-%d-%d%s", time.Now().UnixNano(), rand.Int63(), index, path.Ext(file_name))
	}

	return new_file_name_list
}

//创建博客功能点静态资源存放的目录
//返回目录
func (s *AttachmentBaseService) createBlogDir() string {
	dir := strings.Join([]string{UPLOAD_ROOT_PATH, "blog"}, "/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	return dir
}
