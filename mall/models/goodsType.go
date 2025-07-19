package models

type GoodsType struct {
	Id          int
	Title       string
	Description string
	Status      int
	AddTime     int64
	UpdateTime  int64
}

func (GoodsType) TableName() string {
	return "goods_type"
}
