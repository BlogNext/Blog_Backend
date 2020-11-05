package entity

type BaseEntity struct {
	DocID     string `json:"doc_id"`     //文档在es中的唯一标识和es的_index是一样的，这里赋值一下
	ID        uint64 `json:"id"`         //文档在数据库中的唯一标识
	CreatedAt uint64 `json:"created_at"` //文档在数据库的创建时间
	UpdatedAt uint64 `json:"updated_at"` //文档在数据库的更新时间
}
