package models

type Photo struct {
	GormModel
	Title    string    `gorm:"not null" json:"title" form:"title"`
	Caption  string    `json:"caption,omitempty" form:"caption"`
	PhotoUrl string    `gorm:"not null" json:"photo_url" form:"photo_url"`
	Comments []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
	UserId   uint
}
