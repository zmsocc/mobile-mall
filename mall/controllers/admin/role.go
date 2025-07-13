package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

type RoleController struct {
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
		Title:       title,
		Description: description,
		Status:      1,
		AddTime:     time.Now().Unix(),
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

func (c RoleController) Auth(ctx *gin.Context) {
	// 获取角色 id
	roleId, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/role")
		return
	}
	// 获取所有的权限
	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Preload("AccessItem").Find(&accessList)
	
	// 获取当前角色拥有的权限，并把权限 id 放在一个 map 对象里面
	roleAccess := []models.RoleAccess{}
	models.DB.Where("role_id=?", roleId).Find(&roleAccess)
	roleAccessMap := make(map[int]int)
	for _, v := range roleAccess {
		roleAccessMap[v.AccessId] = v.AccessId
	}

	// 循环遍历所有的权限数据，判断当前的权限的 id 是否在角色权限的 Map 对象中，如果是的话，给当前数据加入 checked 属性
	for i := 0; i < len(accessList); i++ {
		if _, ok := roleAccessMap[accessList[i].Id]; ok {
			accessList[i].Checked = true
		}
		for j := 0; j < len(accessList[i].AccessItem); j++ {
			if _, ok := roleAccessMap[accessList[i].AccessItem[j].Id]; ok {
			accessList[i].AccessItem[j].Checked = true
			}
		}
	}

	ctx.HTML(http.StatusOK, "admin/role/auth.html", gin.H{
		"roleId":     roleId,
		"accessList": accessList,
	})
}

func (c RoleController) DoAuth(ctx *gin.Context) {
	roleId, err := models.Int(ctx.PostForm("role_id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/role")
		return
	}
	// 获取权限 id 切片
	accessIds := ctx.PostFormArray("access_node[]")
	// 删除当前角色对应的权限

	// 增加数据
	roleAccess := models.RoleAccess{}
	for _, v := range accessIds {
		roleAccess.RoleId = roleId
		accessId, _ := models.Int(v)
		roleAccess.AccessId = accessId
		models.DB.Create(&roleAccess)
	}
	ctx.String(http.StatusOK, "DoAUth")
}
