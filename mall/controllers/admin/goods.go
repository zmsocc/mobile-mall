package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

type GoodsController struct {
	BaseController
}

func (c GoodsController) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/goods/index.html", gin.H{})
}

func (c GoodsController) Add(ctx *gin.Context) {
	// 获取商品分类
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = ?", 0).Preload("GoodsCateItems").Find(&goodsCateList)
	
	// 获取所有颜色信息
	goodsColorlist := []models.GoodsColor{}
	models.DB.Find(&goodsColorlist)

	// 获取商品规格包装
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	ctx.HTML(http.StatusOK, "admin/goods/add.html", gin.H{
		"goodsCateList": goodsCateList,
		"goodsColorlist": goodsColorlist,
		"goodsTypeList": goodsTypeList,
	})
}

func (c GoodsController) DoAdd(ctx *gin.Context) {
	attrIdList := ctx.PostFormArray("attr_id_list")
	attrValueList := ctx.PostFormArray("attr_value_list")
	goodsImageList := ctx.PostFormArray("goods_image_list")
	ctx.JSON(http.StatusOK, gin.H{
		"attrIdList": attrIdList,
		"attrValueList": attrValueList, 
		"goodsImageList": goodsImageList,
	})
}

func (c GoodsController) ImageUpload(ctx *gin.Context) {
	// 上传图片
	// 可以在网络里面看到传递的参数
	imageDir, err := models.UploadImg(ctx, "file")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"link": "",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"link": "/"+imageDir,
	})
}

func (c GoodsController) GoodsTypeAttribute(ctx *gin.Context) {
	cateId, err := models.Int(ctx.Query("cateId"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false, 
			"result": "",
		})
		return
	}
	goodsTypeAttributeList := []models.GoodsTypeAttribute{}
	err = models.DB.Where("cate_id = ?", cateId).Find(&goodsTypeAttributeList).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false, 
			"result": "",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true, 
		"result": goodsTypeAttributeList,
	})
}