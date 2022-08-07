package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 使用内存池减少对象的频繁申请和回收

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.Run()
}
