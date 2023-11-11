package responses

type AnalysisScoreResult struct {
	Id      string `json:"id"`
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}
