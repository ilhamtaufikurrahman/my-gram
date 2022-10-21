package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type SocialMedia struct {
	GormModel
	Name           string `gorm:"not null" json:"name" form:"name" valid:"required~Your name is required"`
	SocialMediaUrl string `gorm:"not null" json:"social_media_url" form:"social_media_url" valid:"required~Your social media url is required"`
	User           *User  `json:"user"`
	UserId         uint   `json:"user_id"`
}

type UserSocialMedia struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetSocialMedias struct {
	Id             uint            `json:"id"`
	Name           string          `json:"name"`
	SocialMediaUrl string          `json:"social_media_url"`
	UserId         uint            `json:"user_id"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	User           UserSocialMedia `json:"user"`
}

func (sc *SocialMedia) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(sc)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

func (sc *SocialMedia) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(sc)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}
