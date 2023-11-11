package responses

type ErrorResponse struct {
	Err     string `json:"error_code"`
	ErrCode int    `json:"err_code"`
}
