package db

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IlyaKhalizov/golang-united-certificates/pkg/models"

	"github.com/google/uuid"
)

type InMemDb struct {
	records []models.Certificate
}

func (rcv *InMemDb) GetCertById(certId string) (result models.Certificate, err error) {
	for _, cert := range rcv.records {
		if cert.CertId == certId {
			return cert, err
		}
	}
	return result, errors.New("No cert was found")
}

func (rcv *InMemDb) IsCertExistsByUserAndCourse(userId, courseId string) bool {
	for _, cert := range rcv.records {
		if cert.UserId == userId && cert.CourseId == courseId {
			return true
		}
	}
	return false
}

func (rcv *InMemDb) Issue(userId, courseId string) (result models.Certificate, err error) {
	result = models.Certificate{CertId: uuid.New().String(), UserId: userId, CourseId: courseId, CreatedAt: time.Now()}
	rcv.records = append(rcv.records, result)
	return
}

func (rcv *InMemDb) List(pageSize int, pageToken string) (result []models.Certificate, NextPageToken string, err error) {
	if pageSize == 0 {
		return rcv.records, NextPageToken, nil
	}
	pt, err := strconv.Atoi(pageToken)
	if pt >= len(rcv.records) {
		return nil, NextPageToken, errors.New("Incorrect page token")
	}
	if pt+pageSize >= len(rcv.records) {
		result = append(result, rcv.records[pt:len(rcv.records)]...)
	} else {
		result = append(result, rcv.records[pt:pt+pageSize]...)
		NextPageToken = strconv.Itoa(pt + pageSize)
	}
	return result, NextPageToken, nil
}

func (rcv *InMemDb) ListForUser(pageSize int, pageToken string, userId string) (result []models.Certificate, NextPageToken string, err error) {
	intermediateDb := InMemDb{}
	for _, cert := range rcv.records {
		if cert.UserId == userId {
			intermediateDb.records = append(intermediateDb.records, cert)
		}
	}
	if len(intermediateDb.records) == 0 {
		return result, NextPageToken, errors.New("No certificates were found for user")
	}
	return intermediateDb.List(pageSize, pageToken)
}

func (rcv *InMemDb) ListForCourse(pageSize int, pageToken string, courseId string) (result []models.Certificate, NextPageToken string, err error) {
	intermediateDb := InMemDb{}
	for _, cert := range rcv.records {
		if cert.UserId == courseId {
			intermediateDb.records = append(intermediateDb.records, cert)
		}
	}
	if len(intermediateDb.records) == 0 {
		return []models.Certificate{}, NextPageToken, errors.New("No certificates were found for course")
	}
	return intermediateDb.List(pageSize, pageToken)
}

func (rcv *InMemDb) Delete(certId string) {
	for k, cert := range rcv.records {
		if cert.CertId == certId {
			rcv.records[k], rcv.records[len(rcv.records)-1] = rcv.records[len(rcv.records)-1], rcv.records[k]
		}
	}
	rcv.records = rcv.records[0 : len(rcv.records)-2]
}

func (rcv *InMemDb) Connect(connectionString string) {
	log.Println("initilazing local In-Memory Database...")
	rcv.init()
	log.Println("done!")
}

// func newIMD generates data for in-memory storage
func (rcv *InMemDb) init() {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Duration(2))
		rcv.records = append(rcv.records, models.Certificate{CertId: fmt.Sprint(uuid.New()), UserId: fmt.Sprint(uuid.New()), CourseId: fmt.Sprint(uuid.New()), CreatedAt: time.Now()})
	}
}
