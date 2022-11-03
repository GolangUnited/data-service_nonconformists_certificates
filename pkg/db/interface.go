// Package db contains DB interface
// and it's implementations for different backends.

package db

import "github.com/koldyba/gu-certmgr/pkg/models"

type DB interface {
	Connect(connectionString string)
	GetCertById(certId string) (result models.Certificate, err error)
	IsCertExistsByUserAndCourse(userId, courseId string) bool
	Issue(userId, courseId string) (result models.Certificate, err error)
	List(pageSize int, pageToken string) (result []models.Certificate, NextPageToken string, err error)
	ListForUser(pageSize int, pageToken string, userId string) (result []models.Certificate, NextPageToken string, err error)
	ListForCourse(pageSize int, pageToken string, courseId string) (result []models.Certificate, NextPageToken string, err error)
	Delete(certId string)
}
