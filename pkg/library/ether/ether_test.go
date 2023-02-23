package ether

import (
	"fmt"
	"testing"

	"github.com/yenole/chainx/pkg/model"
)

var (
	btst = &model.Chain{CID: 97, URL: "https://bsctst.showmeta.io"}
)

func TestBlockNumber(t *testing.T) {
	fmt.Println(BlockNumber(btst))
}

func TestCall(t *testing.T) {
	byts, err := Call(btst, []byte(eBlockNumber))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("string(byts): %v\n", string(byts))
}
