package models

type Focus struct {
	Id         int
	Title      string
	FocusType  int
	FocusImg   string
	Link       string
	Sort       int
	Status     int
	AddTime    int64
	UpdateTime int64
}

func (Focus) TableName() string {
	return "focus"
}
