package admin

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
	"gorm.io/gorm"
)

type MainController struct{}

func (con MainController) Index(ctx *gin.Context) {
	//获取userinfo 对应的session
	session := sessions.Default(ctx)
	userinfo := session.Get("userinfo")
	//类型断言 来判断 userinfo是不是一个string
	userinfoStr, ok := userinfo.(string)

	if ok {
		// 获取用户信息
		var userinfoStruct []models.Manager
		json.Unmarshal([]byte(userinfoStr), &userinfoStruct)
		// 	fmt.Println(userinfoStruct)
	  	// 获取所有的权限 
		accessList := []models.Access{}
		models.DB.Where("module_id=?", 0).Preload("AccessItem", func(db *gorm.DB) *gorm.DB {
			return db.Order("access.sort DESC")
		}).Order("sort DESC").Find(&accessList)
		
		// 获取当前角色拥有的权限，并把权限 id 放在一个 map 对象里面
		roleAccess := []models.RoleAccess{}
		models.DB.Where("role_id=?", userinfoStruct[0].RoleId).Find(&roleAccess)
		roleAccessMap := make(map[int]int)
		for _, v := range roleAccess {
			roleAccessMap[v.AccessId] = v.AccessId
		}

		// 循环遍历所有的权限数据，判断当前权限的 id 是否在角色权限的 Map 对象中，如果是的话给当前数据加入 checked 属性
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

		ctx.HTML(http.StatusOK, "admin/main/index.html",gin.H{
			"username": userinfoStruct[0].Username,
			"accessList": accessList,
			"isSuper": userinfoStruct[0].IsSuper,
		})
	} else {
		ctx.Redirect(302, "/admin/login")
	}
}

func (con MainController) Welcome(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/main/welcome.html", gin.H{})
}

func (con MainController) ChangeStatus(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false, 
			"message": "传入的参数错误",
		})
		return
	} 
	table := ctx.Query("table")
	field := ctx.Query("field")
	var currentStatus int
	models.DB.Table(table).Where("id=?", id).Pluck(field, &currentStatus)
	newStatus := 1
	if currentStatus == 1 {
		newStatus = 2
	}
	err = models.DB.Table(table).Where("id=?", id).Updates(map[string]interface{}{
		field: newStatus,
	}).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false, 
			"message": "更新失败"+err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"status": newStatus,
	})
}

func (c MainController) ChangeNum(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false, 
			"message": "传入的参数错误",
		})
		return
	}
	table := ctx.Query("table")
	field := ctx.Query("field")
	num := ctx.Query("num")
	err = models.DB.Exec("update "+table+" set "+field+"="+num+" where id = ?", id).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false, 
			"message": "跟新数据失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true, 
		"message": "跟新数据成功",
	})
}
