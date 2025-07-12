package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/controllers/occ"
)

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", occ.DefaultController{}.Index)
		defaultRouters.GET("/news", occ.DefaultController{}.News)

	}
}
