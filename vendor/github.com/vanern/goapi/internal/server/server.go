package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	Address string
	Gin     *gin.Engine
}

var Gins []Server

func RegisterGin(address string) *gin.Engine {
	gn := gin.New()

	Gins = append(Gins, Server{
		Address: address,
		Gin:     gn,
	})
	return gn
}

// assume you have something like:
//
//	var endpoints = []typs.Endpoint{ ... }
func Start() {

	for _, srvr := range Gins {
		srvr.Gin.Run(srvr.Address)
	}

}
