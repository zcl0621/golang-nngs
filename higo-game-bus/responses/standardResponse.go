package responses

type StandardResponse struct {
	Code int         `json:"code"`    // 0 成功 1失败
	Msg  string      `json:"message"` // 提示语
	Data interface{} `json:"data"`
}
