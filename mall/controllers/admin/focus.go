package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
	"net/http"
)

type FocusController struct {
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
		Title:     title,
		FocusType: focusType,
		FocusImg:  focusImg,
		Link:      link,
		Sort:      sort,
		Status:    status,
		AddTime:   models.GetUnix(),
	}
	err = models.DB.Create(&focus).Error
	if err != nil {
		c.Error(ctx, "增加轮播图失败", "/admin/focus/add")
		return
	}
	c.Success(ctx, "增加轮播图成功", "/admin/focus")
}

func (c FocusController) Edit(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "获取数据失败", "/admin/focus/edit")
		return
	}
	focus := models.Focus{Id: id}
	models.DB.Find(&focus)
	ctx.HTML(http.StatusOK, "admin/focus/edit.html", gin.H{
		"focus": focus,
	})
}

func (c FocusController) DoEdit(ctx *gin.Context) {
	id, _ := models.Int(ctx.Query("id"))
	title := ctx.PostForm("title")
	focusType, err := models.Int(ctx.PostForm("focus_type"))
	if err != nil {
		c.Error(ctx, "获取focus_type失败", "/admin/focus/edit?id="+models.String(id))
		return
	}
	link := ctx.PostForm("link")
	sort, err := models.Int(ctx.PostForm("sort"))
	if err != nil {
		c.Error(ctx, "获取sort失败", "/admin/focus/edit?id="+models.String(id))
		return
	}
	status, err := models.Int(ctx.PostForm("status"))
	if err != nil {
		c.Error(ctx, "获取status失败", "/admin/focus/edit?id="+models.String(id))
		return
	}
	// 上传图片
	focusImg, _ := models.UploadImg(ctx, "focus_img")

	focus := models.Focus{Id: id}
	focus.Id = id
	focus.Title = title
	focus.FocusType = focusType
	focus.Link = link
	focus.Sort = sort
	focus.Status = status
	focus.UpdateTime = models.GetUnix()
	if focus.FocusImg != "" {
		focus.FocusImg = focusImg	
	}
	if err := models.DB.Save(&focus).Error; err != nil {
		c.Error(ctx, "修改失败", "/admin/focus/edit?id="+models.String(id))
		return
	}
	c.Success(ctx, "修改成功", "/admin/focus")
}

func (c FocusController) Delete(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "获取数据失败", "/admin/focus")
		return
	}
	focus := models.Focus{Id: id}
	models.DB.Delete(&focus)
	// 根据自己的需要，看是软删除还是要真删除
	// oc.Remove("static/upload/20250718/1752848755.jpg")
	c.Success(ctx, "删除数据成功", "/admin/focus")
}
