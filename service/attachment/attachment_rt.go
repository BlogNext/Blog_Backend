package attachment

import (
	"github.com/blog_backend/entity/attachment"
	"github.com/blog_backend/model"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type AttachmentRtService struct {
	AttachmentBaseService
}

//网络下载博客功能点的静态资源
//url 网络下载的资源
//module 功能点
//file_type 文件类型
func (s *AttachmentRtService) DownloadBlogImage(url string, module, file_type int64) (full_attachment_extend []*attachment.AttachmentEntity) {
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
