package occ

import (
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	. "github.com/hunterhug/go_image"
	qrcode "github.com/skip2/go-qrcode"
)

type DefaultController struct{}

func (c DefaultController) Index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "首页")
}

func (c DefaultController) ThumbnailV1(ctx *gin.Context) {
	// 按宽度进行比例缩放，输入输出都是文件
	// filename string, savepath string, width int
	filename := "static/upload/0.png"
	savepath := "static/upload/0_600.png"
	err := ScaleF2F(filename, savepath, 600)
	if err != nil {
		ctx.String(http.StatusOK, "生成图片失败")
		return
	}
	ctx.String(http.StatusOK, "ThumbnailV1 生成图片成功")
}

func (c DefaultController) ThumbnailV2(ctx *gin.Context) {
	filename := "static/upload/0.png"
	savepath := "static/upload/0_400.png"
	// 按宽、高进行比例缩放，输入输出都是文件
	// 如果宽高比例和以前不一样的话，就会执行剪切操作
	// filename string, savepath string, width int, height int
	err := ThumbnailF2F(filename, savepath, 400, 200)
	if err != nil {
		ctx.String(http.StatusOK, "生成图片失败")
		return
	}
	ctx.String(http.StatusOK, "ThumbnailV2 生成图片成功")
}

func (c DefaultController) QrcodeV1(ctx *gin.Context) {
	var png []byte
	png, err := qrcode.Encode("https://www.occ.com", qrcode.Medium, 256)
	if err != nil {
		ctx.String(http.StatusOK, "生成二维码失败")
		return
	}
	ctx.String(http.StatusOK, string(png))
}

func (c DefaultController) QrcodeV2(ctx *gin.Context) {
	savepath := "static/upload/qrcode.png"
	err := qrcode.WriteFile("https://www.occ.com", qrcode.Medium, 256, savepath)
	if err != nil {
		ctx.String(http.StatusOK, "生成二维码失败")
		return
	}
	file, _ := os.ReadFile(savepath)
	ctx.String(http.StatusOK, string(file))
}


// func (con DefaultController) Index(c *gin.Context) {
// 	//设置sessions
// 	session := sessions.Default(c)
// 	//配置session的过期时间
// 	session.Options(sessions.Options{
// 		MaxAge: 3600 * 6, // 6hrs   MaxAge单位是秒
// 	})
// 	session.Set("username", "张三 111")
// 	session.Save() //设置session的时候必须调用

// 	c.HTML(http.StatusOK, "default/index.html", gin.H{
// 		"msg": "我是一个msg",
// 		"t":   models.GetUnix(),
// 	})
// }
// func (con DefaultController) News(c *gin.Context) {
// 	//获取sessions
// 	session := sessions.Default(c)
// 	username := session.Get("username")
// 	c.String(200, "username=%v", username)
// }
