package internal

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx *gin.Context) interface{}

type wrap struct {
	*gin.RouterGroup
}

func Wrap(r *gin.Engine) *wrap {
	return &wrap{RouterGroup: &r.RouterGroup}
}

func (r *wrap) Group(relativePath string, handlers ...gin.HandlerFunc) *wrap {
	return &wrap{r.RouterGroup.Group(relativePath, handlers...)}
}

func (r *wrap) GET(relativePath string, handler HandlerFunc, handlers ...gin.HandlerFunc) *wrap {
	r.RouterGroup.GET(relativePath, append(handlers, r.handle(handler))...)
	return r
}

func (r *wrap) POST(relativePath string, handler HandlerFunc) *wrap {
	r.RouterGroup.POST(relativePath, r.handle(handler))
	return r
}

func (r *wrap) handle(fun HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var buff [4096]byte
				n := runtime.Stack(buff[:], false)
				fmt.Printf("==> %s\n", string(buff[:n]))
				ctx.String(http.StatusInternalServerError, fmt.Sprint("err:", err))
			}
		}()
		ret := fun(ctx)
		switch v := ret.(type) {
		case error:
			ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": v.Error()})

		case func():
			v()

		default:
			result := gin.H{"code": 200}
			if v != nil {
				result["data"] = v
			}
			ctx.JSON(http.StatusOK, result)
		}

	}
}
