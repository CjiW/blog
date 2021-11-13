package routers

import "github.com/gin-gonic/gin"

func UseRouter(r *gin.Engine) {
	r.POST("/")
}
