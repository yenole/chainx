package api

import (
	"strconv"
)

func str2uint(str string) uint {
	num, _ := strconv.Atoi(str)
	return uint(num)
}

func rsperr(message string) *Result {
	return &Result{
		Id:      "null",
		Jsonrpc: "2.0",
		Error: &RspError{
			Code:    -200,
			Message: message,
		},
	}
}
