package models

type GoodsImage struct {
	Id         int
	GoodsId    int
	ImgUrl     string
	ColorId    int
	Sort       int
	Status     int
	AddTime    int64
	UpdateTime int64
}

func (GoodsImage) TableName() string {
	return "goods_image"
}
