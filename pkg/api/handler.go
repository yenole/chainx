package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yenole/chainx/pkg/api/internal"
	"github.com/yenole/chainx/pkg/library/ether"
)

func (s *Service) setupRouting() {
	api := internal.Wrap(s.rt)
	api = api.Group("v1")

	api.POST(`/:chain`, s.handleChain)
}

func (s *Service) handleChain(c *gin.Context) interface{} {
	chain := str2uint(c.Param("chain"))
	var body json.RawMessage
	c.BindJSON(&body)
	ch := s.safeChain(chain)
	if ch == nil {
		return func() { c.JSON(http.StatusOK, rsperr("chain not support")) }
	}

	count := 0
try:
	count += 1
	raw, err := ether.Call(ch, body)
	if err != nil {
		if count <= 3 {
			goto try
		}
		return func() { c.JSON(http.StatusOK, rsperr("chain rpc fail")) }
	}
	return func() { c.JSON(http.StatusOK, json.RawMessage(raw)) }
}
