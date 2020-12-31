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
func (s *AttachmentBkService) UploadBlog(Ctx *gin.Context) (fullAttachmentExtend []*attachment.AttachmentEntity) {
	multipartForm, _ := Ctx.MultipartForm()
	files := multipartForm.File["upload_blog_images[]"]
	modules := multipartForm.Value["modules[]"]
	fileTypes := multipartForm.Value["file_type[]"]

	//创建博客功能点静态资源存放的目录
	dir := s.createBlogDir()

	var attachmentIds []uint64

	//保存成功的文件
	attachmentIds = make([]uint64, len(files))

	for index, file := range files {

		//获取重命名的文件名
		renameList := s.renameFileName([]string{file.Filename})
		fileRename := renameList[0]

		dst := strings.Join([]string{dir, fileRename}, "/")

		//验证modules和file_types是否合法
		attachmentModel := new(model.AttachmentModel)

		var module int64
		var fileType int64

		module, err := strconv.ParseInt(modules[index], 10, 64)
		if err != nil {
			panic(exception.NewException(exception.VALIDATE_ERR, "modules参数无法转化成整形"))
		}

		fileType, err = strconv.ParseInt(fileTypes[index], 10, 64)

		if err != nil {
			panic(exception.NewException(exception.VALIDATE_ERR, "file_type参数无法转化成整形"))
		}

		attachmentModel.CheckValidModule(module)
		attachmentModel.CheckValidFileType(fileType)

		//上传文件
		err = Ctx.SaveUploadedFile(file, dst)
		if err != nil {
			panic(err)
		}

		//保存到数据库
		attachmentModel = s.saveToDB(dst, module, fileType)

		attachmentIds = append(attachmentIds, uint64(attachmentModel.ID))
	}

	//获取文件列表
	fullAttachmentExtend = GetAttachmentImages(attachmentIds)

	return
}
