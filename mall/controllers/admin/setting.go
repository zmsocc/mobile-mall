package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

type SettingController struct {
	BaseController
}

func (c SettingController) Index(ctx *gin.Context) {
	setting := models.Setting{}
	models.DB.First(&setting)
	ctx.HTML(http.StatusOK, "admin/setting/index.html", gin.H{
		"setting": setting,
	})
}

func (c SettingController) DoEdit(ctx *gin.Context) {
	now := models.GetUnix()
	setting := models.Setting{Id: 1}
	models.DB.Find(&setting)
	if err := ctx.ShouldBind(&setting); err != nil {
		c.Error(ctx, "绑定数据失败", "/admin/setting")
		return 
	}
	siteLogo, err := models.UploadImg(ctx, "site_logo")
	if err == nil && len(siteLogo) > 0 {
		setting.SiteLogo = siteLogo
	}
	noPicture, err := models.UploadImg(ctx, "no_picture")
	if err == nil && len(noPicture) > 0 {
		setting.NoPicture = noPicture
	}
	setting.UpdateTime = now
	err = models.DB.Save(&setting).Error
	if err != nil {
		c.Error(ctx, "修改数据失败", "/admin/setting")
		return 
	}
	c.Success(ctx, "修改数据成功", "/admin/setting")
}