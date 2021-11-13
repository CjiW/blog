package main

import (
	"Blog/middleware"
	"Blog/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(middleware.Cors())
	routers.UseRouter(r)
	_ = r.Run(":8081")
}
