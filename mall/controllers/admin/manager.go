package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

type ManagerController struct{
	BaseController
}

func (c ManagerController) Index(ctx *gin.Context) {
	managerList := []models.Manager{}
	models.DB.Preload("Role").Find(&managerList)
	// fmt.Printf("%#v", managerList)
	ctx.HTML(http.StatusOK, "admin/manager/index.html", gin.H{
		"managerList": managerList,
	})
}

func (c ManagerController) Add(ctx *gin.Context) {
	roleList := []models.Role{}
	models.DB.Find(&roleList)
	ctx.HTML(http.StatusOK, "admin/manager/add.html", gin.H{
		"roleList": roleList,
	})
}

func (c ManagerController) DoAdd(ctx *gin.Context) {
	roleId, err := models.Int(ctx.PostForm("role_id"))
	if err != nil {
		c.Error(ctx, "获取数据失败", "/admin/manager/add")
		return
	}
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	mobile := ctx.PostForm("mobile")
	email := ctx.PostForm("email")
	if len(username) < 2 || len(password) < 6 {
		c.Error(ctx, "用户名或者密码的长度不合法", "/admin/manager/add")
		return
	}
	// 判断管理员是否存在
	managerList := []models.Manager{}
	models.DB.Where("username=?", username).Find(&managerList)
	if len(managerList) > 0 {
		c.Error(ctx, "此管理员已存在", "/admin/manager/add")
		return
	}
	// 执行增加管理员
	managerInfo := models.Manager{
		Username: username,
		Password: models.Md5(password),
		Mobile: mobile,
		Email: email,
		Status: 1,
		RoleId: roleId,
		AddTime: time.Now().Unix(),
	}
	err = models.DB.Create(&managerInfo).Error
	if err != nil {
		c.Error(ctx, "增加管理员失败", "/admin/manager/add")
		return
	}
	c.Success(ctx, "增加管理员成功", "/admin/manager")
}

func (c ManagerController) Edit(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "获取id失败", "/admin/manager")
		return
	}
	managerInfo := models.Manager{Id: id}
	models.DB.Find(&managerInfo)
	roleList := []models.Role{}
	models.DB.Find(&roleList)
	ctx.HTML(http.StatusOK, "admin/manager/edit.html", gin.H{
		"managerInfo": managerInfo,
		"roleList": roleList,
	})
}

func (c ManagerController) DoEdit(ctx *gin.Context) {
	id, err := models.Int(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "获取信息失败", "/admin/manager/edit")
		return
	}
	roleId, err := models.Int(ctx.PostForm("role_id"))
	if err != nil {
		c.Error(ctx, "获取信息失败", "/admin/manager/edit")
		return
	}
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	email := ctx.PostForm("email")
	mobile := ctx.PostForm("mobile")
	if len(username) < 2 {
		c.Error(ctx, "用户名长度不合法", "/admin/manager/edit?id="+models.String(id))
		return
	}
	if password != "" {
		if len(password) < 6 {
			c.Error(ctx, "密码的长度不合法", "/admin/manager/edit?id="+models.String(id))
			return
		}
	}
	if len(mobile) > 11 {
		c.Error(ctx, "手机号的长度不合法", "/admin/manager/edit?id="+models.String(id))
		return
	}
	managerInfo := models.Manager{Id: id}
	models.DB.Find(&managerInfo)
	managerInfo.Username = username
	managerInfo.Password = models.Md5(password)
	managerInfo.Email = email
	managerInfo.Mobile = mobile
	managerInfo.RoleId = roleId
	managerInfo.AddTime = time.Now().Unix()
	err = models.DB.Save(&managerInfo).Error
	if err != nil {
		c.Error(ctx, "修改失败", "/admin/manager/edit")
		return
	}
	c.Success(ctx, "修改成功", "/admin/manager")
}

func (c ManagerController) Delete(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "获取信息失败", "/admin/manager/edit")
		return
	}
	managerInfo := models.Manager{Id: id}
	models.DB.Delete(&managerInfo)
	c.Success(ctx, "删除成功", "/admin/manager")
}
