package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/mobile-mall/mall/models"
	"net/http"
)

type GoodsCateController struct {
	BaseController
}

func (c GoodsCateController) Index(ctx *gin.Context) {
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = ?", 0).Preload("GoodsCateItems").Find(&goodsCateList)
	// fmt.Printf("%#v", goodsCateList)
	ctx.HTML(http.StatusOK, "admin/goodsCate/index.html", gin.H{
		"goodsCateList": goodsCateList,
	})
}

func (c GoodsCateController) Add(ctx *gin.Context) {
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = ?", 0).Find(&goodsCateList)
	ctx.HTML(http.StatusOK, "admin/goodsCate/add.html", gin.H{
		"goodsCateList": goodsCateList,
	})
}

func (c GoodsCateController) DoAdd(ctx *gin.Context) {
	title := ctx.PostForm("title")
	pid, err := models.Int(ctx.PostForm("pid"))
	if err != nil {
		c.Error(ctx, "传入参数错误", "/goodsCate/add")
		return
	}
	link := ctx.PostForm("link")
	template := ctx.PostForm("template")
	subTitle := ctx.PostForm("sub_title")
	keywords := ctx.PostForm("keywords")
	description := ctx.PostForm("description")
	sort, err := models.Int(ctx.PostForm("sort"))
	if err != nil {
		c.Error(ctx, "传入sort参数错误", "/admin/goodsCate/add")
		return
	}
	status, err := models.Int(ctx.PostForm("status"))
	if err != nil {
		c.Error(ctx, "传入status参数错误", "/admin/goodsCate/add")
		return
	}
	cateImgDir, _ := models.UploadImg(ctx, "cate_img")
	now := models.GetUnix()
	goodsCate := models.GoodsCate{
		Title:       title,
		Pid:         pid,
		SubTitle:    subTitle,
		Link:        link,
		Template:    template,
		Keywords:    keywords,
		Description: description,
		CateImg:     cateImgDir,
		Sort:        sort,
		Status:      status,
		AddTime:     now,
		UpdateTime:  now,
	}
	err = models.DB.Create(&goodsCate).Error
	if err != nil {
		c.Error(ctx, "增加商品失败", "/admin/goodsCate/add")
		return
	}
	c.Success(ctx, "增加商品成功", "/admin/goodsCate")
}

func (c GoodsCateController) Edit(ctx *gin.Context) {
	// 获取要修改的数据
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/goodsCate")
		return
	}
	goodsCate := models.GoodsCate{Id: id}
	models.DB.Find(&goodsCate)

	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = ?", 0).Preload("GoodsCateItems").Find(&goodsCateList)
	ctx.HTML(http.StatusOK, "admin/goodsCate/index.html", gin.H{
		"goodsCateList": goodsCateList,
	})

	ctx.HTML(http.StatusOK, "admin/goodsCate/edit.html", gin.H{
		"goodsCate":     goodsCate,
		"goodsCateList": goodsCateList,
	})
}

func (c GoodsCateController) DoEdit(ctx *gin.Context) {
	id, err := models.Int(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "传入参数错误", "/goodsCate/add")
		return
	}
	title := ctx.PostForm("title")
	pid, err := models.Int(ctx.PostForm("pid"))
	if err != nil {
		c.Error(ctx, "传入参数错误", "/goodsCate/add")
		return
	}
	link := ctx.PostForm("link")
	template := ctx.PostForm("template")
	subTitle := ctx.PostForm("sub_title")
	keywords := ctx.PostForm("keywords")
	description := ctx.PostForm("description")
	sort, err := models.Int(ctx.PostForm("sort"))
	if err != nil {
		c.Error(ctx, "传入sort参数错误", "/admin/goodsCate/edit?id="+models.String(id))
		return
	}
	status, err := models.Int(ctx.PostForm("status"))
	if err != nil {
		c.Error(ctx, "传入status参数错误", "/admin/goodsCate/edit?id="+models.String(id))
		return
	}
	cateImgDir, _ := models.UploadImg(ctx, "cate_img")
	now := models.GetUnix()
	goodsCate := models.GoodsCate{Id: id}
	models.DB.Find(&goodsCate)
	goodsCate.Title = title
	goodsCate.Pid = pid
	goodsCate.Link = link
	goodsCate.Template = template
	goodsCate.SubTitle = subTitle
	goodsCate.Keywords = keywords
	goodsCate.Description = description
	goodsCate.Sort = sort
	goodsCate.Status = status
	goodsCate.UpdateTime = now
	if cateImgDir != "" {
		goodsCate.CateImg = cateImgDir
	}
	if err := models.DB.Save(&goodsCate).Error; err != nil {
		c.Error(ctx, "更新商品数据失败", "/admin/goodsCate/edit?id="+models.String(id))
		return
	}
	c.Success(ctx, "更新商品数据成功", "/admin/goodsCate")
}

func (c GoodsCateController) Delete(ctx *gin.Context) {
	id, err := models.Int(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "传入数据错误", "/admin/goodsCate")
		return
	}
	// 获取要删除的数据
	goodsCate := models.GoodsCate{Id: id}
	models.DB.Find(&goodsCate)
	if goodsCate.Pid != 0 {
		models.DB.Delete(&goodsCate)
		c.Success(ctx, "删除商品成功", "/admin/goodsCate")
		return
	}
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = ?", goodsCate.Id).Find(&goodsCateList)
	if len(goodsCateList) > 0 {
		c.Error(ctx, "当前分类下面含有子分类，请删除子分类后再来删除这个数据", "/admin/goodsCate")
		return
	}
	models.DB.Delete(&goodsCate)
	c.Success(ctx, "删除商品成功", "/admin/goodsCate")
}
