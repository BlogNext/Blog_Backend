package es


type BaseDoc struct {
	ID         uint64 `json:"id"`     //文档在数据库中的唯一标识
	CreateTime uint64 `json:"create_time"`  //文档在数据库的创建时间
	UpdateTime uint64 `json:"update_time"`  //文档在数据库的更新时间
}
