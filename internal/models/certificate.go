package models

import "time"

type Certificate struct {
	Id        string
	UserId    string
	CourseId  string
	CreatedAt time.Time
}
