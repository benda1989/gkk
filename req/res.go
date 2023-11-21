package req

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

type List struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}
