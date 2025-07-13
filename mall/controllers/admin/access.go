package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

type AccessController struct{
	BaseController
}

func (c AccessController) Index(ctx *gin.Context) {
	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Preload("AccessItem").Find(&accessList)
	// fmt.Println(AccessList)
	ctx.HTML(http.StatusOK, "admin/access/index.html", gin.H{
		"accessList": accessList,
	})
}

func (c AccessController) Add(ctx *gin.Context) {
	// 获取顶级模块
	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Find(&accessList)
	ctx.HTML(http.StatusOK, "admin/access/add.html", gin.H{
		"accessList": accessList,
	})
}

func (c AccessController) DoAdd(ctx *gin.Context) {
	moduleName := strings.Trim(ctx.PostForm("module_name"), " ")
	accessType, err1 := models.Int(ctx.PostForm("type"))
	actionName := ctx.PostForm("action_name")
	url := ctx.PostForm("url")
	moduleId, err2 := models.Int(ctx.PostForm("module_id"))
	sort, err3 := models.Int(ctx.PostForm("sort"))
	status, err4 := models.Int(ctx.PostForm("status"))
	description := ctx.PostForm("description")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.Error(ctx, "获取数据失败", "/admin/access/add")
		return
	}
	if moduleName == "" {
		c.Error(ctx, "模块名称不能为空", "/admin/access/add")
		return
	}
	access := models.Access{
		ModuleName: moduleName,
		ActionName: actionName,
		Type: accessType,
		Url: url,
		ModuleId: moduleId,
		Sort: sort,
		Status: status,
		Description: description,
		AddTime: time.Now().Unix(),
	}
	err := models.DB.Create(&access).Error
	if err != nil {
		c.Error(ctx, "增加权限管理失败", "/admin/access/add")
		return
	}
	c.Success(ctx, "增加权限管理成功", "/admin/access")
}

func (c AccessController) Edit(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "获取数据失败", "/admin/access/edit")
		return
	}
	access := models.Access{Id: id}
	models.DB.Preload("AccessItem").Find(&access)
	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Find(&accessList)
	ctx.HTML(http.StatusOK, "admin/access/edit.html", gin.H{
		"access": access,
		"accessList": accessList,
	})
}

func (c AccessController) DoEdit(ctx *gin.Context) {
	id, err := models.Int(ctx.PostForm("id"))
	moduleName := strings.Trim(ctx.PostForm("module_name"), " ")
	accessType, err1 := models.Int(ctx.PostForm("type"))
	actionName := ctx.PostForm("action_name")
	url := ctx.PostForm("url")
	moduleId, err2 := models.Int(ctx.PostForm("module_id"))
	sort, err3 := models.Int(ctx.PostForm("sort"))
	status, err4 := models.Int(ctx.PostForm("status"))
	description := ctx.PostForm("description")
	if err != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.Error(ctx, "获取数据失败", "/admin/access/edit?id="+models.String(id))
		return
	}
	if moduleName == "" {
		c.Error(ctx, "模块名称不能为空", "/admin/access/edit?id="+models.String(id))
		return
	}
	access := models.Access{Id: id}
	access.ModuleName = moduleName
	access.ActionName = actionName
	access.Type = accessType
	access.Url = url
	access.ModuleId = moduleId
	access.Sort = sort
	access.Status = status
	access.Description = description
	access.AddTime = time.Now().Unix()
	
	err = models.DB.Save(&access).Error
	if err != nil {
		c.Error(ctx, "增加权限管理失败", "/admin/access/add")
		return
	}
	c.Success(ctx, "增加权限管理成功", "/admin/access")
}

func (c AccessController) Delete(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "查询数据失败", "/admin/access")
		return
	}
	accessInfo := models.Access{Id: id}
	if accessInfo.ModuleId == 0 {	// 表示它是顶级模块
		mId := accessInfo.Id
		access := []models.Access{}
		models.DB.Where("module_id=?", mId).Find(&access)
		if len(access) > 0 {
			c.Error(ctx, "当前模块下还有内容，删除失败", "/admin/access")
			return
		}
	}
	models.DB.Delete(&accessInfo)
	c.Success(ctx, "删除成功", "/admin/access")
}
