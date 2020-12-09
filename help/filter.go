package help

//过滤的工具和entity.ListResponseEntity配合使用
type Option struct {
	Label string      `json:"label"` //中文提示
	Value interface{} `json:"value"` //具体是值，具体的类型
}

type Filter struct {
	Label   string   `json:"label"` //中文提示
	Field   string   `json:"field"` //传输字段的名字
	Options []Option `json:"options"`
}
