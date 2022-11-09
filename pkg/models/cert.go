package models

import "time"

type Certificate struct {
	Id string
	// Id        string `gorm:"index:idx_uuid_pagination"`
	UserId    string
	CourseId  string
	CreatedAt time.Time
	// CreatedAt time.Time `gorm:"index:idx_uuid_pagination"`
}
