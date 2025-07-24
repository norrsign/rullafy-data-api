package framework

import (
	"github.com/gin-gonic/gin"
	"github.com/vanern/goapi/internal/server"
)

func RegisterGin(address string) *gin.Engine {
	return server.RegisterGin(address)
}
