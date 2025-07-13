package models

type Access struct {
	Id          int
	ModuleName  string // 模块名称
	Type        int    // 节点类型：1.表示模块 2.表示菜单 3.表示操作
	ActionName  string // 操作名称
	Url         string // 路由跳转地址
	ModuleId    int    // 此 module_id 和当前模型的_id关联，module_id=0 表示模块
	Sort        int
	Description string
	Status      int
	AddTime     int64
	AccessItem  []Access `gorm:"foreignKey:ModuleId;references:Id"`
	Checked     bool     `gorm:"-"` // 忽略本字段
}

func (Access) TableName() string {
	return "access"
}
