package models

type SocialMedia struct {
	GORMModel
	Name           string `gorm:"not null" json:"name" valid:"required~name is required"`
	SocialMediaURL string `gorm:"not null" json:"social_media_url" valid:"required~social_media_url is required"`
	UserID         uint   `gorm:"not null" json:"user_id" valid:"required~user_id is required"`
	User           *User  `json:"user"`
}
