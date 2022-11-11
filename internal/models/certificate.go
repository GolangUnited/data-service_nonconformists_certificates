package models

import "time"

// Certificate is base model of object to operate
// within golang-united-certificates
type Certificate struct {
	Id        string
	UserId    string
	CourseId  string
	CreatedAt time.Time
	CreatedBy string
	DeletedAt time.Time
	DeletedBy string
}
