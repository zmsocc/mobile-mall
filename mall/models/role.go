package models

type Role struct {
	Id int
	Title string
	Description string
	Status int
	AddTime int64
}

func (Role) TableName() string {
	return "role"
}