package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})
	server.Run(":8080")
}
