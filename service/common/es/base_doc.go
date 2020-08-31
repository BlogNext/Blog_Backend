package es

const (
	//博客索引
	BLOG_INDEX = "xiaochen_blog_next_blog"
)

type BaseDoc struct {
	DocID      string `json:"doc_id"`      //文档在es中的唯一标识和es的_index是一样的，这里赋值一下
	ID         uint64 `json:"id"`          //文档在数据库中的唯一标识
	CreateTime uint64 `json:"create_time"` //文档在数据库的创建时间
	UpdateTime uint64 `json:"update_time"` //文档在数据库的更新时间
}
