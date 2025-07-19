package occ

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

type DefaultController struct{}

func (con DefaultController) Index(c *gin.Context) {
	//设置sessions
	session := sessions.Default(c)
	//配置session的过期时间
	session.Options(sessions.Options{
		MaxAge: 3600 * 6, // 6hrs   MaxAge单位是秒
	})
	session.Set("username", "张三 111")
	session.Save() //设置session的时候必须调用

	c.HTML(http.StatusOK, "default/index.html", gin.H{
		"msg": "我是一个msg",
		"t":   models.GetUnix(),
	})
}
func (con DefaultController) News(c *gin.Context) {
	//获取sessions
	session := sessions.Default(c)
	username := session.Get("username")
	c.String(200, "username=%v", username)
}
