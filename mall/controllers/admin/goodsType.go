package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
	"net/http"
	"strings"
)

type GoodsTypeController struct {
	BaseController
}

func (c GoodsTypeController) Index(ctx *gin.Context) {
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)
	ctx.HTML(http.StatusOK, "admin/goodsType/index.html", gin.H{
		"goodsTypeList": goodsTypeList,
	})
}

func (c GoodsTypeController) Add(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/goodsType/add.html", gin.H{})
}

func (c GoodsTypeController) DoAdd(ctx *gin.Context) {
	// 从前端获取数据并去除空格
	title := strings.Trim(ctx.PostForm("title"), " ")
	description := strings.Trim(ctx.PostForm("description"), " ")
	if title == "" {
		c.Error(ctx, "标题不能为空", "/admin/goodsType/add")
		return
	}
	status, err := models.Int(ctx.PostForm("status"))
	if err != nil {
		c.Error(ctx, "传入status错误", "/admin/goodsType/add")
		return
	}
	now := models.GetUnix()
	goodsTypeList := models.GoodsType{
		Title:       title,
		Description: description,
		Status:      status,
		AddTime:     now,
		UpdateTime:  now,
	}
	err = models.DB.Create(&goodsTypeList).Error
	if err != nil {
		c.Error(ctx, "增加商品类型失败，请重试", "/admin/goodsType/add")
		return
	}
	c.Success(ctx, "增加商品类型成功", "/admin/goodsType")
}

func (c GoodsTypeController) Edit(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/goodsType")
		return
	}
	goodsType := models.GoodsType{Id: id}
	models.DB.Find(&goodsType)
	ctx.HTML(http.StatusOK, "admin/goodsType/edit.html", gin.H{
		"goodsType": goodsType,
	})
	c.Success(ctx, "更新成功", "/admin/goodsType/edit")
}

func (c GoodsTypeController) DoEdit(ctx *gin.Context) {
	id, err := models.Int(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/goodsType/edit?id="+models.String(id))
		return
	}
	title := strings.Trim(ctx.PostForm("title"), " ")
	description := strings.Trim(ctx.PostForm("description"), " ")
	status, err := models.Int(ctx.PostForm("status"))
	if err != nil {
		c.Error(ctx, "传入status错误", "/admin/goodsType/edit?id="+models.String(id))
		return
	}
	if title == "" {
		c.Error(ctx, "标题不能为空", "/admin/goodsType/edit?id="+models.String(id))
		return
	}
	goodsType := models.GoodsType{Id: id}
	models.DB.Find(&goodsType)
	goodsType.Title = title
	goodsType.Description = description
	goodsType.Status = status
	goodsType.UpdateTime = models.GetUnix()
	err = models.DB.Save(&goodsType).Error
	if err != nil {
		c.Error(ctx, "保存失败", "/admin/goodsType/edit?id="+models.String(id))
		return
	}
	c.Success(ctx, "修改数据成功", "/admin/goodsType")
}

func (c GoodsTypeController) Delete(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/goodsType")
		return
	}
	goodsType := models.GoodsType{Id: id}
	models.DB.Delete(&goodsType)
	c.Success(ctx, "成功删除数据", "/admin/goodsType")
}
