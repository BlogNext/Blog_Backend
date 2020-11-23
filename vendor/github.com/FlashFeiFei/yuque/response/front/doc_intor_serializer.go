package front

type ResponseDocIntorSerializer struct {
	Data *DocIntorSerializer `json:"data"`
}

/**
文档简介（调用前端接口的,非开放平台的api，以后可能会出问题）
*/
type DocIntorSerializer struct {
	ID                int64  `json:"id"`                 //博客id
	BookId            int64  `json:"book_id"`            //知识库id
	Cover             string `json:"cover"`              //博客封面图
	CustomDescription string `json:"custom_description"` //博客摘要
}
