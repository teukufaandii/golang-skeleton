package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Name     string    `gorm:"size:100;not null" json:"name"`
	Email    string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Phone    string    `gorm:"size:20" json:"phone"`
	Password string    `gorm:"size:255;not null" json:"-"`
	RoleID   uuid.UUID `gorm:"type:uuid;not null;index" json:"role_id"`

	Role Role `gorm:"foreignKey:RoleID" json:"role"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(ctx *gorm.DB) error {
	err := u.BaseModel.BeforeCreate(ctx)
	if err != nil {
		return err
	}

	if u.RoleID == uuid.Nil {
		var role Role
		if err := ctx.Where("code = ?", "user").First(&role).Error; err != nil {
			return err
		}
		u.RoleID = role.ID
	}

	return nil
}

func (u *User) IsAdmin() bool {
	return u.Role.Code == "admin"
}
