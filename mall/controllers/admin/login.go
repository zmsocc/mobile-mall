package admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

type LoginController struct {
	BaseController
}

func (con LoginController) Index(ctx *gin.Context) {
	//验证md5是否正确
	// fmt.Println(models.Md5("123456"))   e10adc3949ba59abbe56e057f20f883e
	
	ctx.HTML(http.StatusOK, "admin/login/login.html", gin.H{})
}
func (con LoginController) DoLogin(ctx *gin.Context) {

	captchaId := ctx.PostForm("captchaId")
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	verifyValue := ctx.PostForm("verifyValue")
	// fmt.Println(username, password)
	//1、验证验证码是否正确
	if flag := models.VerifyCaptcha(captchaId, verifyValue); flag {
		//2、查询数据库 判断用户以及密码是否存在
		userinfoList := []models.Manager{}
		password = models.Md5(password)

		models.DB.Where("username=? AND password=?", username, password).Find(&userinfoList)

		if len(userinfoList) > 0 {
			//3、执行登录 保存用户信息 执行跳转
			session := sessions.Default(ctx)
			//注意：session.Set没法直接保存结构体对应的切片 把结构体转换成json字符串
			userinfoSlice, _ := json.Marshal(userinfoList)
			session.Set("userinfo", string(userinfoSlice))
			session.Save()
			con.Success(ctx, "登录成功", "/admin")

		} else {
			con.Error(ctx, "用户名或者密码错误", "/admin/login")
		}

	} else {
		con.Error(ctx, "验证码验证失败", "/admin/login")
	}

}

func (con LoginController) Captcha(ctx *gin.Context) {
	id, b64s, err := models.MakeCaptcha()

	if err != nil {
		fmt.Println(err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
}
func (con LoginController) LoginOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("userinfo")
	session.Save()
	con.Success(ctx, "退出登录成功", "/admin/login")
}
