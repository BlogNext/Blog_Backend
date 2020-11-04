package attachment

import (
	"fmt"
	"github.com/blog_backend/common-lib/config"
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/entity/attachment"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
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

	if len(attachment_list) <= 0 {
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
		attachment_entity.FileType = attachment_model.FileType
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
func (s *AttachmentService) saveToDB(dst string, module, file_type int64) (attachment_model *model.AttachmentModel) {

	db := mysql.GetDefaultDBConnect()
	attachment_model = new(model.AttachmentModel)
	attachment_model.Path = dst
	attachment_model.CreateTime = time.Now().Unix()
	attachment_model.UpdateTime = time.Now().Unix()
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
func (s *AttachmentService) renameFileName(file_name_list []string) (new_file_name_list []string) {

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
func (s *AttachmentService) createBlogDir() string {
	dir := strings.Join([]string{UPLOAD_ROOT_PATH, "blog"}, "/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	return dir
}

//网络下载博客功能点的静态资源
//url 网络下载的资源
//module 功能点
//file_type 文件类型
func (s *AttachmentService) DownloadBlogImage(url string, module, file_type int64) (full_attachment_extend []*attachment.AttachmentEntity) {
	response, err := http.Get(url)

	if err != nil {
		return nil
	}

	//创建博客功能点静态资源存放的目录
	dir := s.createBlogDir()

	//获取重命名的文件名
	rename_list := s.renameFileName([]string{url})
	file_rename := rename_list[0]

	dst := strings.Join([]string{dir, file_rename}, "/")

	//验证modules和file_types是否合法
	attachment_model := new(model.AttachmentModel)
	attachment_model.CheckValidModule(module)
	attachment_model.CheckValidFileType(file_type)

	log.Println(dst)
	//创建一个文件
	out, err := os.Create(dst)
	if err != nil {
		return nil
	}

	defer out.Close()

	_, err = io.Copy(out, response.Body)

	if err != nil {
		panic(err)
	}
	log.Println("完美拷贝")

	//数据库保存
	attachment_model = s.saveToDB(dst, module, file_type)

	//返回数据
	full_attachment_extend = GetAttachmentImages([]uint64{uint64(attachment_model.ID)})

	return full_attachment_extend
}

//上传博客功能点的文件
func (s *AttachmentService) UploadBlog(Ctx *gin.Context) (full_attachment_extend []*attachment.AttachmentEntity) {
	multipart_form, _ := Ctx.MultipartForm()
	files := multipart_form.File["upload_blog_images[]"]
	modules := multipart_form.Value["modules[]"]
	file_types := multipart_form.Value["file_type[]"]

	//创建博客功能点静态资源存放的目录
	dir := s.createBlogDir()

	var attachment_ids []uint64

	//保存成功的文件
	attachment_ids = make([]uint64, len(files))

	for index, file := range files {

		//获取重命名的文件名
		rename_list := s.renameFileName([]string{file.Filename})
		file_rename := rename_list[0]

		dst := strings.Join([]string{dir, file_rename}, "/")

		//验证modules和file_types是否合法
		attachment_model := new(model.AttachmentModel)

		var module int64
		var file_type int64

		module, err := strconv.ParseInt(modules[index], 10, 64)
		if err != nil {
			panic(exception.NewException(exception.VALIDATE_ERR, "modules参数无法转化成整形"))
		}

		file_type, err = strconv.ParseInt(file_types[index], 10, 64)

		if err != nil {
			panic(exception.NewException(exception.VALIDATE_ERR, "file_type参数无法转化成整形"))
		}

		attachment_model.CheckValidModule(module)
		attachment_model.CheckValidFileType(file_type)

		//上传文件
		err = Ctx.SaveUploadedFile(file, dst)
		if err != nil {
			panic(err)
		}

		//保存到数据库
		attachment_model = s.saveToDB(dst, module, file_type)

		attachment_ids = append(attachment_ids, uint64(attachment_model.ID))
	}

	//获取文件列表
	full_attachment_extend = GetAttachmentImages(attachment_ids)

	return
}
