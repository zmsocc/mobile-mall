package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
	"gopkg.in/ini.v1"
)

func InitAdminAuthMiddleware(ctx *gin.Context) {
	// excludeAuthPath("aaa")
	//进行权限判断 没有登录的用户 不能进入后台管理中心
	// fmt.Println("InitAdminAuthMiddleware")
	// 获取 url 访问的地址
	// /admin/captcha?t=0.643991812512453
	pathName := strings.Split(ctx.Request.URL.String(), "?")[0]
	// fmt.Println(pathName)
	
	// 获取 session 里面保存的用户信息
	session := sessions.Default(ctx)
	userinfo := session.Get("userinfo")
	// 类型断言，来判断 userinfo 是不是一个 string
	userinfoStr, ok := userinfo.(string)
	if !ok {
		// 用户没有登陆
		if pathName != "/admin/login" && pathName != "/admin/doLogin" && pathName != "/admin/captcha" {
			ctx.Redirect(302, "/admin/login")
		}
		return
	}
	var userinfoStruct []models.Manager
	err := json.Unmarshal([]byte(userinfoStr), &userinfoStruct)
	if err != nil || len(userinfoStruct) <= 0 || userinfoStruct[0].Username == "" {
		if pathName != "/admin/login" && pathName != "/admin/doLogin" && pathName != "/admin/captcha" {
			ctx.Redirect(302, "/admin/login")
		}
		return
	} else { 	// 用户登陆成功，权限判断
		urlPath := strings.Replace(pathName, "/admin/", "", 1)
		if userinfoStruct[0].IsSuper != 1 && !excludeAuthPath("/"+urlPath) {
			// 根据角色获取当前角色的权限列表， 然后把权限 id 放在一个 map 类型的对象里面
			roleAccess := []models.RoleAccess{}
			models.DB.Where("role_id=?", userinfoStruct[0].RoleId).Find(&roleAccess)
			roleAccessMap := make(map[int]int)
			for _, v := range roleAccess {
				roleAccessMap[v.AccessId] = v.AccessId
			}

			// 获取当前访问的 url 对应的权限 id， 判断权限 id 是否在角色对应的权限
			// pathname     /admin/manager

			access := models.Access{}
			models.DB.Where("url=?", urlPath).Find(&access)
			if _, ok := roleAccessMap[access.Id]; !ok {
				ctx.String(http.StatusOK, "没有权限")
				ctx.Abort()
			}
		}

	}
}

func excludeAuthPath(urlPath string) bool {
	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		os.Exit(1)
	}
	excludeAuthPath := config.Section("").Key("excludeAuthPath").String()
	excludeAuthPathSlice := strings.Split(excludeAuthPath, ",")
	fmt.Println(excludeAuthPathSlice)

	for _, v := range excludeAuthPathSlice {
		if v == urlPath {
			return true
		}
	}
	return false 
}
