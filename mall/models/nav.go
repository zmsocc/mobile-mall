package models

type Nav struct {
	Id         int
	Title      string
	Link       string
	Position   int
	IsOpennew  int
	Relation   string
	Sort       int
	Status     int
	AddTime    int64
	UpdateTime int64
	GoodsItems []Goods `gorm:"-"` // 忽略本字段
}

func (Nav) TableName() string {
	return "nav"
}
