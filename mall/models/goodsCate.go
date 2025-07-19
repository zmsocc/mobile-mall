package models

type GoodsCate struct {
	Id             int
	Title          string
	Link           string
	CateImg        string
	Template       string
	Pid            int
	SubTitle       string
	Keywords       string
	Description    string
	Status         int
	Sort           int
	AddTime        int64
	UpdateTime     int64
	GoodsCateItems []GoodsCate `gorm:"foreignKey:Pid;references:Id"`
}

func (GoodsCate) TableName() string {
	return "goods_cate"
}
