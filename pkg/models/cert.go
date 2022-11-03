package models

import "time"

type Certificate struct {
	CertId    string
	UserId    string
	CourseId  string
	CreatedAt time.Time
}
