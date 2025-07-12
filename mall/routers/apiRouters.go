package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/controllers/api"
)

func ApiRoutersInit(r *gin.Engine) {
	apiRouters := r.Group("/api")
	{
		apiRouters.GET("/", api.ApiController{}.Index)
		apiRouters.GET("/userlist", api.ApiController{}.Userlist)
		apiRouters.GET("/plist", api.ApiController{}.Plist)
	}

}
