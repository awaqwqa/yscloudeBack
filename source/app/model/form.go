package model

type LoadForm struct {
	ServerNum int           `json:"server_num"`
	LoadPos   [3]int        `json:"pos"`
	IsAdvance bool          `json:"advance"`
	Ac        AdvanceConfig `json:"config"`
}

// 用于返回信息的
type BackMsg struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Body   interface{} `json:"body"`
}
type RegisterForm struct {
	RedeemKey string `json:"redeem_key" binding:"required"`
	UserName  string `json:"user_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	QQ        int    `json:"QQ" binding:"required"`
}

// login form
type LoginForm struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdvanceConfig struct {
	BlockState    bool `json:"block_state"`
	ClearOldBlock bool `json:"clear_old_block"`
	RebootBack    bool `json:"no_back"`
	ClearItems    bool `json:"clear_items"`
	AddStruct     bool `json:"add_struct"`
}
