package admin

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

type FocusController struct{
	BaseController
}

func (c FocusController) Index(ctx *gin.Context) {
	focusList := []models.Focus{}
	models.DB.Find(&focusList)
	ctx.HTML(http.StatusOK, "admin/focus/index.html", gin.H{
		"focusList": focusList,
	})
}
func (c FocusController) Add(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/focus/add.html", gin.H{})
}

func (c FocusController) DoAdd(ctx *gin.Context) {
	title := ctx.PostForm("title")
	focusType, err := models.Int(ctx.PostForm("focus_type"))
	if err != nil {
		c.Error(ctx, "非法请求", "/admin/focus/add")
		return
	}
	link := ctx.PostForm("link")
	sort, err := models.Int(ctx.PostForm("sort"))
	if err != nil {
		c.Error(ctx, "请输入正确的排序值", "/admin/focus/add")
		return
	}
	status, err := models.Int(ctx.PostForm("status"))
	if err != nil {
		c.Error(ctx, "非法请求", "/admin/focus/add")
		return
	}

	// 上传图片
	focusImg, err := models.UploadImg(ctx, "focus_img")
	if err != nil {
		fmt.Println(err)
	}
	focus := models.Focus{
		Title: title,
		FocusType: focusType,
		FocusImg: focusImg,
		Link: link,
		Sort: sort,
		Status: status,
		AddTime: models.GetUnix(),
	}
	err = models.DB.Create(&focus).Error
	if err != nil {
		c.Error(ctx, "增加轮播图失败", "/admin/focus/add")
	}
	c.Success(ctx, "增加轮播图成功", "/admin/focus")
}

func (c FocusController) Edit(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/focus/edit.html", gin.H{})
}

func (c FocusController) DoEdit(ctx *gin.Context) {
	ctx.String(http.StatusOK, "DoEdit")
}

func (c FocusController) Delete(ctx *gin.Context) {
	ctx.String(http.StatusOK, "-add--文章-")
}
