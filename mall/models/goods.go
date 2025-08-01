package models

type Goods struct {
	Id            int
	Title         string
	SubTitle      string
	GoodsSn       string
	CateId        int
	ClickCount    int
	GoodsNumber   int
	Price         float64
	MarketPrice   float64
	RelationGoods string
	GoodsAttr     string
	GoodsColor    string
	GoodsVersion  string
	GoodsImg      string
	GoodsGift     string
	GoodsFitting  string
	GoodsKeywords string
	GoodsDesc     string
	GoodsContent  string
	IsDelete      int
	IsHot         int
	IsBest        int
	IsNew         int
	GoodsTypeId   int
	Sort          int
	Status        int
	AddTime       int64
	UpdateTime    int64
}

func (Goods) TableName() string {
	return "goods"
}

/*
	根据商品分类获取推荐商品
	@param {Number} cateId - 分类id
	@param {String} goodsType - hot best new
	@param {Number} limitNum - 数量
*/

func GetGoodsByCategory(cateId int, goodsType string, limitNum int) []Goods {
	// 判断 cateId 是否是顶级分类
	goodsCate := GoodsCate{Id: cateId}
	DB.Find(&goodsCate)
	tempSlice := []int{}
	// 顶级分类
	if goodsCate.Pid == 0 {
		// 获取顶级分类下面的二级分类
		goodsCateList := []GoodsCate{}
		DB.Where("pid = ?", goodsCate.Id).Find(&goodsCateList)
		for i := 0; i < len(goodsCateList); i++ {
			tempSlice = append(tempSlice, goodsCateList[i].Id)
		}
	}
	tempSlice = append(tempSlice, cateId)

	goodsList := []Goods{}
	where := "cate_id = ?"
	switch goodsType {
	case "hot":
		where += " AND is_hot=1 AND is_delete!=1"
		break
	case "best":
		where += " AND is_best=1 AND is_delete!=1"
		break
	case "new":
		where += " AND is_new=1 AND is_delete!=1"
		break
	default:
		break
	}
	DB.Where(where, cateId).Limit(limitNum).Order("sort DESC").Find(&goodsList)
	return goodsList
}
