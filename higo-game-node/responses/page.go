package responses

type PageResponse struct {
	Count   int64       `json:"count"`   //总数
	Results interface{} `json:"results"` //数据
}
