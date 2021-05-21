package attachment

import (
	"errors"
	"fmt"
	"github.com/blog_backend/common-lib/config"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity/attachment"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"gorm.io/gorm"
	"log"
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
func GetAttachmentImagesMap(ids []uint64) (attachmentEntityList map[uint]*attachment.AttachmentEntity) {
	attachmentList := getAttachmentByIds(ids)
	if attachmentList == nil {
		return
	}

	if len(attachmentList) <= 0 {
		return
	}

	attachmentEntityList = make(map[uint]*attachment.AttachmentEntity, len(attachmentList))

	serverConfig, _ := config.GetConfig("server")
	serverInfo := serverConfig.GetStringMap("server")
	fileDomain := serverInfo["file_domain"].(string)

	for _, attachmentModel := range attachmentList {
		attachmentEntity := new(attachment.AttachmentEntity)
		attachmentEntity.ID = uint64(attachmentModel.ID)
		attachmentEntity.CreatedAt = uint64(attachmentModel.CreatedAt)
		attachmentEntity.UpdatedAt = uint64(attachmentModel.UpdatedAt)
		attachmentEntity.Module = attachmentModel.Module
		attachmentEntity.Path = attachmentModel.Path
		attachmentEntity.Url = attachmentModel.Path
		attachmentEntity.FullUrl = strings.Join([]string{fileDomain, attachmentModel.Path}, "/")
		attachmentEntity.FileType = attachmentModel.FileType

		attachmentEntityList[uint(attachmentEntity.ID)] = attachmentEntity
	}

	return
}

//获取附件
//返回切片
func GetAttachmentImages(ids []uint64) (attachmentEntityList []*attachment.AttachmentEntity) {
	attachmentList := getAttachmentByIds(ids)
	if attachmentList == nil {
		return
	}

	if len(attachmentList) <= 0 {
		return
	}

	attachmentEntityList = make([]*attachment.AttachmentEntity, len(attachmentList))

	serverConfig, _ := config.GetConfig("server")
	serverInfo := serverConfig.GetStringMap("server")
	fileDomain := serverInfo["file_domain"].(string)

	for index, attachmentModel := range attachmentList {
		attachmentEntity := new(attachment.AttachmentEntity)
		attachmentEntity.ID = uint64(attachmentModel.ID)
		attachmentEntity.CreatedAt = uint64(attachmentModel.CreatedAt)
		attachmentEntity.UpdatedAt = uint64(attachmentModel.UpdatedAt)
		attachmentEntity.Module = attachmentModel.Module
		attachmentEntity.Path = attachmentModel.Path
		attachmentEntity.Url = attachmentModel.Path
		attachmentEntity.FullUrl = strings.Join([]string{fileDomain, attachmentModel.Path}, "/")
		attachmentEntity.FileType = attachmentModel.FileType

		attachmentEntityList[index] = attachmentEntity
	}

	return
}

//获取附件
func getAttachmentByIds(ids []uint64) (attachmentList []model.AttachmentModel) {

	if ids == nil {
		return
	}

	db := mysql.GetNewDB(false)
	db.Where("id IN (?)", ids).Find(&attachmentList)

	return
}

//附件服务
type AttachmentBaseService struct {
}

//删除文件
func (s *AttachmentBaseService) DeleteAttachmentById(id uint64){
	db := mysql.GetNewDB(false)
	attachmentModel := new(model.AttachmentModel)
	queryResult := db.Where("id = ?", id).First(attachmentModel)
	notFund := errors.Is(queryResult.Error, gorm.ErrRecordNotFound)
	if notFund {
		return
	}
	//删除文件
	err := os.Remove(attachmentModel.Path)
	if err != nil {
		log.Println("删除文件失败",err)
		return
	}

	db.Delete(attachmentModel)
}

//保存到数据库
func (s *AttachmentBaseService) saveToDB(dst string, module, fileType int64) (attachmentModel *model.AttachmentModel) {

	db := mysql.GetNewDB(false)
	attachmentModel = new(model.AttachmentModel)
	attachmentModel.Path = dst
	attachmentModel.CreatedAt = time.Now().Unix()
	attachmentModel.UpdatedAt = time.Now().Unix()
	attachmentModel.Module = module
	attachmentModel.FileType = fileType

	sqlExecResult := db.Create(attachmentModel)

	if sqlExecResult.Error != nil {
		panic(exception.NewException(exception.DATA_BASE_ERROR_EXEC, fmt.Sprintf("新增失败:%s", sqlExecResult.Error)))
	}

	return attachmentModel
}

//重命名文件名
//file_name_list 一批文件名
func (s *AttachmentBaseService) renameFileName(fileNameList []string) (newFileNameList []string) {

	rand.Seed(time.Now().UnixNano())

	if fileNameList == nil || len(fileNameList) == 0 {
		return nil
	}

	newFileNameList = make([]string, len(fileNameList))
	for index, fileName := range fileNameList {
		newFileNameList[index] = fmt.Sprintf("%d-%d-%d%s", time.Now().UnixNano(), rand.Int63(), index, path.Ext(fileName))
	}

	return newFileNameList
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
