package api

type Result struct {
	Id      string      `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  string      `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type RspError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
