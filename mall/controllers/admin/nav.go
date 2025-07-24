package admin

import (
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

type NavController struct {
	BaseController
}

func (c NavController) Index(ctx *gin.Context) {
	// 获取当前页数
	page, _ := models.Int(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	// 每页数据数量
	pageSize := 5
	// 总共有多少数据
	var count int64
	models.DB.Table("nav").Count(&count)

	navList := []models.Nav{}
	models.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&navList)
	ctx.HTML(http.StatusOK, "admin/nav/index.html", gin.H{
		"navList":    navList,
		"page":       page,
		"totalPages": math.Ceil(float64(count) / float64(pageSize)),
	})
}

func (c NavController) Add(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/nav/add.html", gin.H{})
}

func (c NavController) DoAdd(ctx *gin.Context) {
	// 从前端获取数据并去除空格
	title := strings.Trim(ctx.PostForm("title"), " ")
	if title == "" {
		c.Error(ctx, "角色标题不能为空", "/admin/nav/add")
		return
	}
	link := ctx.PostForm("link")
	position, err := models.Int(ctx.PostForm("position"))
	if err != nil {
		c.Error(ctx, "获取position失败", "/admin/nav/add")
		return
	}
	isOpennew, err := models.Int(ctx.PostForm("is_opennew"))
	if err != nil {
		c.Error(ctx, "获取is_opennew失败", "/admin/nav/add")
		return
	}
	relation := ctx.PostForm("relation")
	sort, err := models.Int(ctx.PostForm("sort"))
	if err != nil {
		c.Error(ctx, "获取sort失败", "/admin/nav/add")
		return
	}
	status, err := models.Int(ctx.PostForm("status"))
	if err != nil {
		c.Error(ctx, "获取status失败", "/admin/nav/add")
		return
	}
	now := models.GetUnix()
	navList := models.Nav{
		Title:      title,
		Link:       link,
		Position:   position,
		IsOpennew:  isOpennew,
		Relation:   relation,
		Sort:       sort,
		Status:     status,
		AddTime:    now,
		UpdateTime: now,
	}
	err = models.DB.Create(&navList).Error
	if err != nil {
		c.Error(ctx, "增加导航失败，请重试", "/admin/nav/add")
		return
	}
	c.Success(ctx, "增加导航成功", "/admin/nav")
}

func (c NavController) Edit(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/nav")
		return
	}
	nav := models.Nav{Id: id}
	models.DB.Find(&nav)
	ctx.HTML(http.StatusOK, "admin/nav/edit.html", gin.H{
		"nav": nav,
	})
}

func (c NavController) DoEdit(ctx *gin.Context) {
	id, err := models.Int(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/nav")
		return
	}
	title := strings.Trim(ctx.PostForm("title"), " ")
	if title == "" {
		c.Error(ctx, "角色的标题不能为空", "/admin/nav/edit")
		return
	}
	link := ctx.PostForm("link")
	position, err := models.Int(ctx.PostForm("position"))
	if err != nil {
		c.Error(ctx, "获取position失败", "/admin/nav/add")
		return
	}
	isOpennew, err := models.Int(ctx.PostForm("is_opennew"))
	if err != nil {
		c.Error(ctx, "获取is_opennew失败", "/admin/nav/add")
		return
	}
	relation := ctx.PostForm("relation")
	sort, err := models.Int(ctx.PostForm("sort"))
	if err != nil {
		c.Error(ctx, "获取sort失败", "/admin/nav/add")
		return
	}
	status, err := models.Int(ctx.PostForm("status"))
	if err != nil {
		c.Error(ctx, "获取status失败", "/admin/nav/add")
		return
	}
	now := models.GetUnix()
	nav := models.Nav{Id: id}
	models.DB.Find(&nav)

	nav.Title = title
	nav.Link = link
	nav.Position = position
	nav.IsOpennew = isOpennew
	nav.Relation = relation
	nav.Sort = sort
	nav.Status = status
	nav.UpdateTime = now

	err = models.DB.Save(&nav).Error
	if err != nil {
		c.Error(ctx, "保存失败", "/admin/nav/edit?id="+models.String(id))
		return
	}
	c.Success(ctx, "修改导航成功", "/admin/nav")
}

func (c NavController) Delete(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/nav")
		return
	}
	nav := models.Nav{Id: id}
	models.DB.Delete(&nav)
	c.Success(ctx, "成功删除数据", "/admin/nav")
}
