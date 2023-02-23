package ether

import (
	"bytes"
	"encoding/json"

	"github.com/yenole/chainx/pkg/library/fetch"
	"github.com/yenole/chainx/pkg/model"
)

const (
	eBlockNumber = `{"jsonrpc":"2.0","id":1,"method":"eth_blockNumber","params":[]}`
)

type result struct {
	Result string `json:"result"`
}

func BlockNumber(c *model.Chain) (uint64, error) {
	var raw result
	err := fetch.Post(c.URL, bytes.NewBufferString(eBlockNumber)).JSON().Apply(&raw).LastErr()
	if err != nil {
		return 0, err
	}
	if v := hex2BigInt(raw.Result); v != nil {
		return v.Uint64(), nil
	}
	return 0, nil
}

func Call(c *model.Chain, body []byte) ([]byte, error) {
	var raw json.RawMessage
	err := fetch.Post(c.URL, bytes.NewBuffer(body)).JSON().Apply(&raw).LastErr()
	if err != nil {
		return nil, err
	}
	return raw, nil
}
