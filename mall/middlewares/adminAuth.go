package middlewares

import (
	"encoding/json"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

func InitAdminAuthMiddleware(ctx *gin.Context) {
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
	}
	// fmt.Println(userinfoStruct)

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"userinfo": userinfoStruct[0].Username,
	// })

}
