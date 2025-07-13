package models

type Manager struct {
	Id       int
	Username string
	Password string
	Mobile   string
	Email    string
	Status   int
	RoleId   int
	AddTime  int64
	IsSuper  int
	Role     Role `gorm:"foreignKey:RoleId"`
}

func (Manager) TableName() string {
	return "manager"
}
