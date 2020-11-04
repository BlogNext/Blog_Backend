package attachment

import (
	"github.com/blog_backend/entity/attachment"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/model"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type AttachmentBkService struct {
	AttachmentBaseService
}


//上传博客功能点的文件
func (s *AttachmentBkService) UploadBlog(Ctx *gin.Context) (full_attachment_extend []*attachment.AttachmentEntity) {
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