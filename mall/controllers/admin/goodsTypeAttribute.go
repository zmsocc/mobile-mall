package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
	"net/http"
	"strings"
)

type GoodsTypeAttributeController struct {
	BaseController
}

func (c GoodsTypeAttributeController) Index(ctx *gin.Context) {
	cateId, err := models.Int(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusOK, "获取id错误")
		return
	}
	// 获取商品类型属性
	goodsTypeAttributeList := []models.GoodsTypeAttribute{}
	models.DB.Where("cate_id = ?", cateId).Find(&goodsTypeAttributeList)
	// 获取商品类型属性对应的类型
	goodsType := models.GoodsType{}
	models.DB.Where("id = ?", cateId).Find(&goodsType)
	ctx.HTML(http.StatusOK, "admin/goodsTypeAttribute/index.html", gin.H{
		"cateId":                 cateId,
		"goodsTypeAttributeList": goodsTypeAttributeList,
		"goodsType":              goodsType,
	})
}

func (c GoodsTypeAttributeController) Add(ctx *gin.Context) {
	// 获取当前商品类型属性对应的类型 id
	cateId, err := models.Int(ctx.Query("cate_id"))
	if err != nil {
		ctx.JSON(http.StatusOK, "获取cate_id错误")
		return
	}
	// 获取所有的商品类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)
	ctx.HTML(http.StatusOK, "admin/goodsTypeAttribute/add.html", gin.H{
		"cateId":        cateId,
		"goodsTypeList": goodsTypeList,
	})
}

func (c GoodsTypeAttributeController) DoAdd(ctx *gin.Context) {
	cateId, err := models.Int(ctx.PostForm("cate_id"))
	if err != nil {
		c.Error(ctx, "传入cateId错误", "/admin/goodsTypeAttribute/add")
		return
	}
	// 从前端获取数据并去除空格
	title := strings.Trim(ctx.PostForm("title"), " ")
	if title == "" {
		c.Error(ctx, "标题不能为空", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
		return
	}
	attrType, err := models.Int(ctx.PostForm("attr_type"))
	if err != nil {
		c.Error(ctx, "传入attrType错误", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
		return
	}
	attrValue := ctx.PostForm("attr_value")
	sort, err := models.Int(ctx.PostForm("sort"))
	if err != nil {
		c.Error(ctx, "传入sort错误", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
		return
	}
	now := models.GetUnix()
	goodsTypeAttributeList := models.GoodsTypeAttribute{
		Title:      title,
		CateId:     cateId,
		AttrType:   attrType,
		AttrValue:  attrValue,
		Sort:       sort,
		Status:     1,
		AddTime:    now,
		UpdateTime: now,
	}
	err = models.DB.Create(&goodsTypeAttributeList).Error
	if err != nil {
		c.Error(ctx, "增加商品属性类型失败，请重试", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
		return
	}
	c.Success(ctx, "增加商品属性类型成功", "/admin/goodsTypeAttribute?id="+models.String(cateId))
}

func (c GoodsTypeAttributeController) Edit(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "获取数据失败", "/admin/goodsTypeAttribute/edit")
		return
	}
	goodsTypeAttribute := models.GoodsTypeAttribute{Id: id}
	models.DB.Find(&goodsTypeAttribute)

	// 获取所有的商品类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	ctx.HTML(http.StatusOK, "admin/goodsTypeAttribute/edit.html", gin.H{
		"goodsTypeAttribute": goodsTypeAttribute,
		"goodsTypeList": goodsTypeList,
	})
}

func (c GoodsTypeAttributeController) DoEdit(ctx *gin.Context) {
	id, err := models.Int(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "获取Id失败", "/admin/goodsType")
		return
	}
	cateId, err := models.Int(ctx.PostForm("cate_id"))
	if err != nil {
		c.Error(ctx, "获取cateId失败", "/admin/goodsType")
		return
	}
	title := strings.Trim(ctx.PostForm("title"), " ")
	if title == "" {
		c.Error(ctx, "标题不能为空", "/admin/goodsTypeAttribute/edit?id="+models.String(id))
		return
	}
	attrType, err := models.Int(ctx.PostForm("attr_type"))
	if err != nil {
		c.Error(ctx, "获取attrType失败", "/admin/goodsTypeAttribute/edit?id="+models.String(id))
		return
	}
	attrValue := ctx.PostForm("attr_value")
	sort, err := models.Int(ctx.PostForm("sort"))
	if err != nil {
		c.Error(ctx, "获取sort失败", "/admin/goodsTypeAttribute/edit?id="+models.String(id))
		return
	}
	now := models.GetUnix()
	// 查找表中数据
	goodsTypeAttribute := models.GoodsTypeAttribute{Id: id}
	models.DB.Find(&goodsTypeAttribute)

	// 执行修改
	goodsTypeAttribute.CateId = cateId
	goodsTypeAttribute.Title = title
	goodsTypeAttribute.AttrType = attrType
	goodsTypeAttribute.AttrValue = attrValue
	goodsTypeAttribute.Sort = sort
	goodsTypeAttribute.UpdateTime = now

	// 改完保存
	if err := models.DB.Save(&goodsTypeAttribute).Error; err != nil {
		c.Error(ctx, "更新商品属性失败", "/admin/goodsTypeAttribute/edit?id="+models.String(id))
		return
	}
	c.Success(ctx, "更新商品属性成功", "/admin/goodsTypeAttribute?id="+models.String(cateId))
}

func (c GoodsTypeAttributeController) Delete(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "获取id失败", "/admin/goodsType")
		return
	}
	cateId, err := models.Int(ctx.Query("cate_id"))
	if err != nil {
		c.Error(ctx, "获取cate_id失败", "/admin/goodsType")
		return
	}
	goodsTypeAttribute := models.GoodsTypeAttribute{Id: id}
	if err := models.DB.Delete(&goodsTypeAttribute).Error; err != nil {
		c.Error(ctx, "删除失败", "/admin/goodsTypeAttribute?id="+models.String(cateId))
		return
	}
	c.Success(ctx, "更新商品属性成功", "/admin/goodsTypeAttribute?id="+models.String(cateId))
}