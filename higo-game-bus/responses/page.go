package responses

type PageResponse struct {
	Count   int64       `json:"count"`   // 总数
	Results interface{} `json:"results"` // 数据
}

type PageRequest struct {
	Page     int `json:"page" form:"page"`           // 页码
	PageSize int `json:"page_size" form:"page_size"` // 每页数量
}
