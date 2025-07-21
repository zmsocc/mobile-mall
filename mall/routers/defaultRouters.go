package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/controllers/occ"
)

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", occ.DefaultController{}.Index)
		defaultRouters.GET("/thumbnailv1", occ.DefaultController{}.ThumbnailV1)
		defaultRouters.GET("/thumbnailv2", occ.DefaultController{}.ThumbnailV2)

		defaultRouters.GET("/qrcodev1", occ.DefaultController{}.QrcodeV1)
		defaultRouters.GET("/qrcodev2", occ.DefaultController{}.QrcodeV2)

	}
}
