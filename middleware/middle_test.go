package middleware

import (
	"fmt"
	"github.com/CodeWizard198/vix"
	"github.com/CodeWizard198/vix/middleware/logging"
	"github.com/CodeWizard198/vix/middleware/recovering"
	"net/http"
	"testing"
)

func TestMiddle(t *testing.T) {
	logMiddle := logging.BuildLogging().Build()
	recoverMiddle := recovering.BuildRecover(http.StatusInternalServerError, []byte("you're panic"))
	recoverMiddle.Operate = func(ctx *vix.Context) {
		fmt.Println("panic... url:", ctx.Req.URL.Path)
	}
	v := vix.NewVIX(vix.ServerWithMiddleware(logMiddle, recoverMiddle.Build()))

	v.GET("/list", func(ctx *vix.Context) {
		ctx.STRING(http.StatusOK, "Hello list, you're successes")
	})

	v.GET("/shop", func(ctx *vix.Context) {
		panic("you're fail, panic...")
		//ctx.STRING(http.StatusOK, "Hello shop, you're successes")
	})

	v.GET("/info", func(ctx *vix.Context) {

	})

	err := v.Run(":8080")
	if err != nil {
		fmt.Println("server fail...")
	}
}
