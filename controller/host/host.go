package host

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Hostname(ctx *gin.Context) {
	host, _ := os.Hostname()

	ctx.String(http.StatusOK, host)

}
