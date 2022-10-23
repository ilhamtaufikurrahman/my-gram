package models

import (
	"my-gram/helpers"

	"gorm.io/gorm"
)

type User struct {
	GormModel
	Email        string        `gorm:"not null;uniqueIndex" json:"email" form:"email"`
	Username     string        `gorm:"not null;uniqueIndex" json:"username" form:"username"`
	Password     string        `gorm:"not null" json:"password" form:"password"`
	Age          uint          `gorm:"not null" json:"age" form:"age"`
	Photos       []Photo       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photos"`
	SocialMedias []SocialMedia `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"social_medias"`
	Comments     []Comment     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password = helpers.HashPass(u.Password)

	return nil
}
