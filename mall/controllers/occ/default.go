package occ

import (
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type DefaultController struct{}

func (c DefaultController) Index(ctx *gin.Context) {
	// 获取顶部导航
	topNavList := []models.Nav{}
	models.DB.Where("status=? AND position=?", 1, 1).Find(&topNavList)
	// 获取轮播图数据
	focusList := []models.Focus{}
	models.DB.Where("status=? AND focus_type=?", 1, 1).Find(&focusList)
	// 获取分类的数据
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=? AND status = ?", 0, 1).Order("sort DESC").Preload("GoodsCateItems", func(db *gorm.DB) *gorm.DB {
		return db.Where("goods_cate.status = ?", 1).Order("goods_cate.sort DESC")
	}).Find(&goodsCateList)
	// 获取中间导航
	middleNavList := []models.Nav{}
	models.DB.Where("status=? AND position=?", 1, 2).Find(&topNavList)

	for i := 0; i < len(middleNavList); i++ {
		relation := strings.ReplaceAll(middleNavList[i].Relation, "，", ",")
		relationIds := strings.Split(relation, ",")
		goodsList := []models.Goods{}
		models.DB.Where("id in ?", relationIds).Select("id, title, goods_img, price").Find(&goodsList)
		middleNavList[i].GoodsItems = goodsList
	}

	// 手机
	phoneList := models.GetGoodsByCategory(1, "best", 8)
	// 配件
	otherList := models.GetGoodsByCategory(2, "all", 1)

	ctx.HTML(http.StatusOK, "itying/index/index.html", gin.H{
		"topNavList":    topNavList,
		"focusList":     focusList,
		"goodsCateList": goodsCateList,
		"middleNavList": middleNavList,
		"phoneList":     phoneList,
		"otherList":     otherList,
	})
}

// func (c DefaultController) ThumbnailV1(ctx *gin.Context) {
// 	// 按宽度进行比例缩放，输入输出都是文件
// 	// filename string, savepath string, width int
// 	filename := "static/upload/0.png"
// 	savepath := "static/upload/0_600.png"
// 	err := ScaleF2F(filename, savepath, 600)
// 	if err != nil {
// 		ctx.String(http.StatusOK, "生成图片失败")
// 		return
// 	}
// 	ctx.String(http.StatusOK, "ThumbnailV1 生成图片成功")
// }

// func (c DefaultController) ThumbnailV2(ctx *gin.Context) {
// 	filename := "C:\\Users\\zhang san\\Desktop\\study\\gogogo\\Project_img\\华为手机05.png"
// 	savepath := "static/upload/huawei_1226X460.png"
// 	// 按宽、高进行比例缩放，输入输出都是文件
// 	// 如果宽高比例和以前不一样的话，就会执行剪切操作
// 	// filename string, savepath string, width int, height int
// 	err := ThumbnailF2F(filename, savepath, 1226, 460)
// 	if err != nil {
// 		ctx.String(http.StatusOK, "生成图片失败")
// 		return
// 	}
// 	ctx.String(http.StatusOK, "ThumbnailV2 生成图片成功")
// }

// func (c DefaultController) QrcodeV1(ctx *gin.Context) {
// 	var png []byte
// 	png, err := qrcode.Encode("https://www.occ.com", qrcode.Medium, 256)
// 	if err != nil {
// 		ctx.String(http.StatusOK, "生成二维码失败")
// 		return
// 	}
// 	ctx.String(http.StatusOK, string(png))
// }

// func (c DefaultController) QrcodeV2(ctx *gin.Context) {
// 	savepath := "static/upload/qrcode.png"
// 	err := qrcode.WriteFile("https://www.occ.com", qrcode.Medium, 256, savepath)
// 	if err != nil {
// 		ctx.String(http.StatusOK, "生成二维码失败")
// 		return
// 	}
// 	file, _ := os.ReadFile(savepath)
// 	ctx.String(http.StatusOK, string(file))
// }

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
