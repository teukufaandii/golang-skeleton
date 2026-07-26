package models

type Role struct {
	BaseModel
	Name string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Code string `gorm:"size:50;uniqueIndex;not null" json:"code"`
}

func (Role) TableName() string {
	return "roles"
}
