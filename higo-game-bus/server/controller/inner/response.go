package inner

type createResponse struct {
	GameId uint `json:"game_id"`
}

type ruleResponse struct {
	SGF              string  `json:"sgf"`
	MaxStep          int     `json:"max_step"`
	TerritoryStep    int     `json:"territory_step"`
	BlackReturnStone float64 `json:"black_return_stone"`
}

type studentHistoryResponse struct {
	LabelName string `json:"label_name"` //对弈类型
	Label     string `json:"-"`
	BlackWin  int    `json:"black_win"`  ///执黑赢数量
	WhiteWin  int    `json:"white_win"`  //执白赢数量
	BlackLose int    `json:"black_lose"` //执黑输数量
	WhiteLose int    `json:"white_lose"` //执白输数量
	Cost      uint   `json:"cost"`       //花费时间 秒
}
