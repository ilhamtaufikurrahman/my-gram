package models

type Comment struct {
	GormModel
	UserId  uint
	PhotoId uint
	Message string `gorm:"not null" json:"message" form:"message"`
}
