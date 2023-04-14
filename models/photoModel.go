package models

// Photo represents the model for an photo
type Photo struct {
	GORMModel
	// gorm.Model
	Title    string    `gorm:"not null" json:"title" valid:"required~title is required"`
	Caption  string    `json:"caption"` // null able
	PhotoURL string    `gorm:"not null" json:"photo_url" valid:"required~photo_url is required"`
	UserID   uint      `gorm:"not null" json:"user_id" valid:"required~user_id is required"`
	User     *User     `json:"user"`
	Comments []Comment `json:"comments"`
}
