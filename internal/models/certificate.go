package models

import "time"

// Certificate is base model of object to operate
// within golang-united-certificates
type Certificate struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid()"`
	UserId    string `gorm:"type:uuid;index:idx_filtering,unique"`
	CourseId  string `gorm:"type:uuid;index:idx_filtering,unique"`
	CreatedAt time.Time
	IsDeleted bool
}
