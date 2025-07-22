package models

type GoodsAttr struct {
	Id              int
	GoodsId         int
	AttributeCateId int
	AttributeId     int
	AttributeTitle  string
	AttributeType   int
	AttributeValue  string
	Sort            int
	Status          int
	AddTime         int64
	UpdateTime      int64
}

func (GoodsAttr) TableName() string {
	return "goods_attr"
}
