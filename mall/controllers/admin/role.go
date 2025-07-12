package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

type RoleController struct{
	BaseController
}

func (c RoleController) Index(ctx *gin.Context) {
	roleList := []models.Role{}
	models.DB.Find(&roleList)
	// fmt.Println(roleList)
	ctx.HTML(http.StatusOK, "admin/role/index.html", gin.H{
		"roleList": roleList,
	})
}

func (c RoleController) Add(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/role/add.html", gin.H{})
}

func (c RoleController) DoAdd(ctx *gin.Context) {
	// 从前端获取数据并去除空格
	title := strings.Trim(ctx.PostForm("title"), " ")
	description := strings.Trim(ctx.PostForm("description"), " ")
	if title == "" {
		c.Error(ctx, "角色标题不能为空", "/admin/role/add")
		return
	}
	roleList := models.Role{
		Title: title,
		Description: description,
		Status: 1,
		AddTime: time.Now().Unix(),
	}
	err := models.DB.Create(&roleList).Error
	if err != nil {
		c.Error(ctx, "增加角色失败，请重试", "/admin/role/add")
		return
	}
	c.Success(ctx, "增加角色成功", "/admin/role")
}

func (c RoleController) Edit(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/role")
		return
	}
	role := models.Role{Id: id}
	models.DB.Find(&role)
	ctx.HTML(http.StatusOK, "admin/role/edit.html", gin.H{
		"role": role,
	})
}

func (c RoleController) DoEdit(ctx *gin.Context) {
	id, err := models.Int(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/role")
		return
	}
	title := strings.Trim(ctx.PostForm("title"), " ")
	description := strings.Trim(ctx.PostForm("description"), " ")
	if title == "" {
		c.Error(ctx, "角色的标题不能为空", "/admin/role/edit")
		return
	}
	role := models.Role{Id: id}
	models.DB.Find(&role)
	role.Title = title
	role.Description = description
	err = models.DB.Save(&role).Error
	if err != nil {
		c.Error(ctx, "保存失败", "/admin/role/edit?id="+models.String(id))
		return
	}
	c.Success(ctx, "修改数据成功", "/admin/role")
}

func (c RoleController) Delete(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/role")
		return
	}
	role := models.Role{Id: id}
	models.DB.Delete(&role)
	c.Success(ctx, "成功删除数据", "/admin/role")
}
