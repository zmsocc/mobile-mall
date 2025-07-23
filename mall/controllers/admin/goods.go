package admin

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
)

var wg sync.WaitGroup

type GoodsController struct {
	BaseController
}

func (c GoodsController) Index(ctx *gin.Context) {
	goodsList := []models.Goods{}
	models.DB.Find(&goodsList)
	ctx.HTML(http.StatusOK, "admin/goods/index.html", gin.H{
		"goodsList": goodsList,
	})
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
		"goodsColorList": goodsColorlist,
		"goodsTypeList": goodsTypeList,
	})
}

func (c GoodsController) DoAdd(ctx *gin.Context) {
	// 获取表单提交过来的数据
	title := ctx.PostForm("title")
	subTitle := ctx.PostForm("sub_title")
	goodsSn := ctx.PostForm("goods_sn")
	cateId, _ := models.Int(ctx.PostForm("cate_id"))
	goodsNumber, _ := models.Int(ctx.PostForm("goods_number"))
	marketPrice, _ := models.Float64(ctx.PostForm("market_price"))
	price, _ := models.Float64(ctx.PostForm("price"))
	relationGoods := ctx.PostForm("relation_goods")
	goodsAttr := ctx.PostForm("goods_attr")
	goodsVersion := ctx.PostForm("goods_version")
	goodsGift := ctx.PostForm("goods_gift")
	goodsFitting := ctx.PostForm("goods_fitting")
	goodsColorArr := ctx.PostFormArray("goods_color")
	goodsKeywords := ctx.PostForm("goods_keywords")
	goodsDesc := ctx.PostForm("goods_desc")
	goodsContent := ctx.PostForm("goods_content")
	isDelete, _ := models.Int(ctx.PostForm("is_delete"))
	isHot, _ := models.Int(ctx.PostForm("is_hot"))
	isBest, _ := models.Int(ctx.PostForm("is_best"))
	isNew, _ := models.Int(ctx.PostForm("is_new"))
	goodsTypeId, _ := models.Int(ctx.PostForm("goods_type_id"))
	sort, _ := models.Int(ctx.PostForm("sort"))
	status, _ := models.Int(ctx.PostForm("status"))
	now := models.GetUnix()

	// 获取颜色信息，把颜色转化成字符串
	goodsColorStr := strings.Join(goodsColorArr, ",")

	// 上传图片，生成缩略图
	goodsImg, _ := models.UploadImg(ctx, "goods_img")

	// 增加商品数据
	goods := models.Goods{
		Title: title,
		SubTitle: subTitle,
		GoodsSn: goodsSn,
		CateId: cateId,
		ClickCount: 100,
		GoodsNumber: goodsNumber,
		MarketPrice: marketPrice,
		Price: price,
		RelationGoods: relationGoods,
		GoodsAttr: goodsAttr,
		GoodsVersion: goodsVersion,
		GoodsGift: goodsGift,
		GoodsFitting: goodsFitting,
		GoodsKeywords: goodsKeywords,
		GoodsDesc: goodsDesc,
		GoodsContent: goodsContent,
		IsDelete: isDelete,
		IsHot: isHot,
		IsBest: isBest,
		IsNew: isNew,
		GoodsTypeId: goodsTypeId,
		Sort: sort,
		Status: status,
		GoodsColor: goodsColorStr,
		GoodsImg: goodsImg,
		AddTime: now,
		UpdateTime: now,
	}
	err := models.DB.Create(&goods).Error
	if err != nil {
		c.Error(ctx, "增加失败", "/admin/goods/add")
		return
	}

	// 增加图库信息
	wg.Add(1)
	go func() {
		goodsImageList := ctx.PostFormArray("goods_image_list")
		for _, v := range goodsImageList {
			goodsImgObj := models.GoodsImage{}
			goodsImgObj.GoodsId = goods.Id
			goodsImgObj.ImgUrl = v
			goodsImgObj.Sort = 10
			goodsImgObj.Status = 1
			goodsImgObj.AddTime = now
			goodsImgObj.UpdateTime = now
			models.DB.Create(&goodsImgObj)
		}
		wg.Done()
	}()
	

	// 增加规格包装
	wg.Add(1)
	go func() {
		attrIdList := ctx.PostFormArray("attr_id_list")
		attrValueList := ctx.PostFormArray("attr_value_list")
		for i := 0; i < len(attrIdList); i++ {
			goodsTypeAttributeId, err := models.Int(attrIdList[i])
			if err != nil {
				ctx.String(http.StatusOK, "转换为整数失败")
				return
			}
			// 获取商品类型属性的数据
			goodsTypeAttributeObj := models.GoodsTypeAttribute{Id: goodsTypeAttributeId}
			models.DB.Find(&goodsTypeAttributeObj)
			// 给商品属性里面增加数据， 规格包装
			goodsAttrObj := models.GoodsAttr{}
			goodsAttrObj.GoodsId = goods.Id
			goodsAttrObj.AttributeTitle = goodsTypeAttributeObj.Title
			goodsAttrObj.AttributeType = goodsTypeAttributeObj.AttrType
			goodsAttrObj.AttributeId = goodsTypeAttributeObj.Id
			goodsAttrObj.AttributeCateId = goodsTypeAttributeObj.CateId
			goodsAttrObj.AttributeValue = attrValueList[i]
			goodsAttrObj.Status = 1
			goodsAttrObj.Sort = 10
			goodsAttrObj.AddTime = now
			goodsAttrObj.UpdateTime = now
			models.DB.Create(&goodsAttrObj)
		}
		wg.Done()
	}()
	wg.Wait()
	c.Success(ctx, "增加商品成功", "/admin/goods")
}

func (c GoodsController) Edit(ctx *gin.Context) {
	// 获取要修改的商品数据
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "传入Id错误", "/admin/goods")
		return
	}
	goods := models.Goods{Id: id}
	models.DB.Find(&goods)

	// 获取商品分类
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = ?", 0).Preload("GoodsCateItems").Find(&goodsCateList)

	// 获取所有颜色信息 
	goodsColorSlice := strings.Split(goods.GoodsColor, ",")
	goodsColorMap := make(map[string]string)
	for _, v := range goodsColorSlice {
		goodsColorMap[v] = v
	}
	goodsColorList := []models.GoodsColor{}
	models.DB.Find(&goodsColorList)
	for i := 0; i < len(goodsColorList); i++ {
		if _, ok := goodsColorMap[models.String(goodsColorList[i].Id)]; ok {
			goodsColorList[i].Checked = true
		}
	}

	// 商品的图库信息
	goodsImageList := []models.GoodsImage{}
	models.DB.Where("goods_id = ?", goods.Id).Find(&goodsImageList)

	// 获取商品类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	// 获取规格信息
	goodsAttr := []models.GoodsAttr{}
	models.DB.Where("goods_id = ?", goods.Id).Find(&goodsAttr)
	goodsAttrStr := ""

	for _, v := range goodsAttr {
		if v.AttributeType == 1 {
			goodsAttrStr += fmt.Sprintf(`<li><span>%v = </span> <input type="hidden" name="attr_id_list" value="%v" /> <input type="text" name="attr_value_list" value="%v" /></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		} else if v.AttributeType == 2 {
			goodsAttrStr += fmt.Sprintf(`<li><span>%v = </span> <input type="hidden" name="attr_id_list" value="%v" /> <input type="text" name="attr_value_list" value="%v" /></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		} else {
			// 获取当前类型对应的值
			goodsTypeAttribute := models.GoodsTypeAttribute{Id: v.AttributeId}
			models.DB.Find(&goodsTypeAttribute)
			attrValueSlice := strings.Split(goodsTypeAttribute.AttrValue, "\n")
			goodsAttrStr += fmt.Sprintf(`<li><span>%v = </span> <input type="hidden" name="attr_id_list" value="%v" />`, v.AttributeTitle, v.AttributeId)
			goodsAttrStr += fmt.Sprintf(`<select name="attr_value_list">`)
			for i := 0; i < len(attrValueSlice); i++ {
				if attrValueSlice[i] == v.AttributeValue {
					goodsAttrStr += fmt.Sprintf(`<option value="%v" selected>%v</option>`, attrValueSlice[i], attrValueSlice[i])
				} else {
					goodsAttrStr += fmt.Sprintf(`<option value="%v">%v</option>`, attrValueSlice)
				}
			}
			goodsAttrStr += fmt.Sprintf(`</select>`)
			goodsAttrStr += fmt.Sprintf(`</li>`)
		}
	}
	// goodsTypeList := []models.GoodsType{}
	// models.DB.Find(&goodsTypeList)

	ctx.HTML(http.StatusOK, "admin/goods/edit.html", gin.H{
		"goods": goods,
		"goodsCateList": goodsCateList, 
		"goodsColorList": goodsColorList,
		"goodsTypeList": goodsTypeList,
		"goodsAttrStr": goodsAttrStr,
		"goodsImageList": goodsImageList,
	})
}

func (c GoodsController) DoEdit(ctx *gin.Context) {
	// 获取表单提交过来的数据
	id, err := models.Int(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "传入id错误", "admin/goods/edit")
	}
	title := ctx.PostForm("title")
	subTitle := ctx.PostForm("sub_title")
	goodsSn := ctx.PostForm("goods_sn")
	cateId, _ := models.Int(ctx.PostForm("cate_id"))
	goodsNumber, _ := models.Int(ctx.PostForm("goods_number"))
	marketPrice, _ := models.Float64(ctx.PostForm("market_price"))
	price, _ := models.Float64(ctx.PostForm("price"))
	relationGoods := ctx.PostForm("relation_goods")
	goodsAttr := ctx.PostForm("goods_attr")
	goodsVersion := ctx.PostForm("goods_version")
	goodsGift := ctx.PostForm("goods_gift")
	goodsFitting := ctx.PostForm("goods_fitting")
	goodsColorArr := ctx.PostFormArray("goods_color")
	goodsKeywords := ctx.PostForm("goods_keywords")
	goodsDesc := ctx.PostForm("goods_desc")
	goodsContent := ctx.PostForm("goods_content")
	isDelete, _ := models.Int(ctx.PostForm("is_delete"))
	isHot, _ := models.Int(ctx.PostForm("is_hot"))
	isBest, _ := models.Int(ctx.PostForm("is_best"))
	isNew, _ := models.Int(ctx.PostForm("is_new"))
	goodsTypeId, _ := models.Int(ctx.PostForm("goods_type_id"))
	sort, _ := models.Int(ctx.PostForm("sort"))
	status, _ := models.Int(ctx.PostForm("status"))
	now := models.GetUnix()

	// 获取颜色信息，把颜色转化成字符串
	goodsColorStr := strings.Join(goodsColorArr, ",")

	goods := models.Goods{Id: id}
	models.DB.Find(&goods)
	goods.Title = title 
	goods.SubTitle = subTitle
	goods.GoodsSn = goodsSn
	goods.CateId = cateId
	goods.ClickCount = 100
	goods.GoodsNumber = goodsNumber
	goods.MarketPrice = marketPrice
	goods.Price = price
	goods.RelationGoods = relationGoods
	goods.GoodsAttr = goodsAttr
	goods.GoodsVersion = goodsVersion
	goods.GoodsGift = goodsGift
	goods.GoodsFitting = goodsFitting
	goods.GoodsKeywords = goodsKeywords
	goods.GoodsDesc = goodsDesc
	goods.GoodsContent = goodsContent
	goods.IsDelete = isDelete
	goods.IsHot = isHot
	goods.IsBest = isBest
	goods.IsNew = isNew
	goods.GoodsTypeId = goodsTypeId
	goods.Sort = sort
	goods.Status = status
	goods.GoodsColor = goodsColorStr
	goods.UpdateTime = now
	// 上传图片，生成缩略图
	goodsImg, err := models.UploadImg(ctx, "goods_img")
	if err == nil && len(goodsImg) > 0 {
		goods.GoodsImg = goodsImg
	}
	
	err = models.DB.Save(&goods).Error
	if err != nil {
		c.Error(ctx, "修改失败", "/admin/goods/edit?id="+models.String(id))
		return
	}

	// 增加图库信息
	wg.Add(1)
	go func() {
		goodsImageList := ctx.PostFormArray("goods_image_list")
		for _, v := range goodsImageList {
			goodsImgObj := models.GoodsImage{}
			goodsImgObj.GoodsId = goods.Id
			goodsImgObj.ImgUrl = v
			goodsImgObj.Sort = 10
			goodsImgObj.Status = 1
			goodsImgObj.AddTime = now
			goodsImgObj.UpdateTime = now
			models.DB.Create(&goodsImgObj)
		}
		wg.Done()
	}()
	

	// 修改规格包装
	// 1.删除当前商品下面的规格包装
	goodsAttrObj := models.GoodsAttr{}
	models.DB.Where("goods_id = ?", goods.Id).Delete(&goodsAttrObj)
	// 2.重新执行增加
	wg.Add(1)
	go func() {
		attrIdList := ctx.PostFormArray("attr_id_list")
		attrValueList := ctx.PostFormArray("attr_value_list")
		for i := 0; i < len(attrIdList); i++ {
			goodsTypeAttributeId, err := models.Int(attrIdList[i])
			if err != nil {
				ctx.String(http.StatusOK, "转换为整数失败")
				return
			}
			// 获取商品类型属性的数据
			goodsTypeAttributeObj := models.GoodsTypeAttribute{Id: goodsTypeAttributeId}
			models.DB.Find(&goodsTypeAttributeObj)
			// 给商品属性里面增加数据，规格包装
			goodsAttrObj := models.GoodsAttr{}
			goodsAttrObj.GoodsId = goods.Id
			goodsAttrObj.AttributeTitle = goodsTypeAttributeObj.Title
			goodsAttrObj.AttributeType = goodsTypeAttributeObj.AttrType
			goodsAttrObj.AttributeId = goodsTypeAttributeObj.Id
			goodsAttrObj.AttributeCateId = goodsTypeAttributeObj.CateId
			goodsAttrObj.AttributeValue = attrValueList[i]
			goodsAttrObj.Status = 1
			goodsAttrObj.Sort = 10
			goodsAttrObj.AddTime = now
			goodsAttrObj.UpdateTime = now
			models.DB.Create(&goodsAttrObj)
		}
		wg.Done()
	}()
	wg.Wait()
	c.Success(ctx, "修改数据成功", "/admin/goods")
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