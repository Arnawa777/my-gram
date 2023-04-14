package models

// Comment represents the model for an comment
type Comment struct {
	GORMModel
	UserID  uint   `gorm:"not null" json:"user_id" valid:"required~user_id is required"`
	PhotoID uint   `gorm:"not null" json:"photo_id" valid:"required~photo_id is required"`
	Message string `gorm:"not null" json:"message" valid:"required~message is required"`
	User    *User  `json:"user"`
	Photo   *Photo `json:"photo"`
}
