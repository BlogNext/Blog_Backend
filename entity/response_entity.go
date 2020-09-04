package entity

import "math"

//通用的http list响应结构体
type ListResponseEntity struct {
	//列表的总数
	Count int64 `json:"count"`
	//总共可以分多少页
	PageCount int `json:"page_count"`
	//前端传过来的参数
	Param interface{} `json:"param,omitempty"`
	//列表的过滤项
	Filter interface{} `json:"filter,omitempty"`
	//列表数据
	List interface{} `json:"list"`

	//额外的一些数据
	Extra interface{} `json:"extra,omitempty"`
}

func (lre *ListResponseEntity) SetCount(count int64) {
	lre.Count = count
}

func (lre *ListResponseEntity) SetPerPage(per_page int) {
	pageCount := float64(lre.Count) / float64(per_page)
	pageCount = math.Ceil(pageCount)
	lre.PageCount = int(pageCount)
}

func (lre *ListResponseEntity) SetParam(param interface{}) {
	lre.Param = param
}

func (lre *ListResponseEntity) SetFilter(filter interface{}) {
	lre.Filter = filter
}

func (lre *ListResponseEntity) SetList(list interface{}) {
	lre.List = list
}

func (lre *ListResponseEntity) SetExtra(extra interface{}) {
	lre.Extra = extra
}
