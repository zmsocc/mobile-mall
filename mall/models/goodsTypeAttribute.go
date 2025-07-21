package models

type GoodsTypeAttribute struct {
	Id         int    `json:"id"`
	CateId     int    `json:"cate_id"`
	Title      string `json:"title"`
	AttrType   int    `json:"attr_type"`
	AttrValue  string `json:"attr_value"`
	Status     int    `json:"status"`
	Sort       int    `json:"sort"`
	AddTime    int64  `json:"add_time"`
	UpdateTime int64  `json:"update_time"`
}

func (GoodsTypeAttribute) TableName() string {
	return "goods_type_attribute"
}
