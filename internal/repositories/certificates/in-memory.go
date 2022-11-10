package db

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"golang-united-certificates/internal/models"

	"github.com/google/uuid"
)

type InMemDb struct {
	records []models.Certificate
}

func (rcv *InMemDb) GetCertById(id string) (models.Certificate, error) {
	for _, cert := range rcv.records {
		if cert.Id == id {
			return cert, nil
		}
	}
	return models.Certificate{}, errors.New("No cert was found")
}

func (rcv *InMemDb) IsCertExistsByUserAndCourse(userId, courseId string) bool {
	for _, cert := range rcv.records {
		if cert.UserId == userId && cert.CourseId == courseId {
			return true
		}
	}
	return false
}

func (rcv *InMemDb) Create(userId, courseId string) (models.Certificate, error) {
	cert := models.Certificate{Id: uuid.New().String(), UserId: userId, CourseId: courseId, CreatedAt: time.Now()}
	rcv.records = append(rcv.records, cert)
	return cert, nil
}

func (rcv *InMemDb) List(pageSize int, pageToken string) ([]models.Certificate, string, error) {
	var npt string
	if pageSize == 0 {
		return rcv.records, npt, nil
	}
	pt, err := strconv.Atoi(pageToken)
	if err != nil {
		return nil, npt, err
	}
	if pt >= len(rcv.records) {
		return nil, npt, errors.New("Incorrect page token")
	}
	var result []models.Certificate
	if pt+pageSize >= len(rcv.records) {
		result = append(result, rcv.records[pt:len(rcv.records)]...)
	} else {
		result = append(result, rcv.records[pt:pt+pageSize]...)
		npt = strconv.Itoa(pt + pageSize)
	}
	return result, npt, nil
}

func (rcv *InMemDb) ListForUser(pageSize int, pageToken string, userId string) ([]models.Certificate, string, error) {
	intermediateDb := InMemDb{}
	for _, cert := range rcv.records {
		if cert.UserId == userId {
			intermediateDb.records = append(intermediateDb.records, cert)
		}
	}
	if len(intermediateDb.records) == 0 {
		return nil, "", errors.New("No certificates were found for user")
	}
	return intermediateDb.List(pageSize, pageToken)
}

func (rcv *InMemDb) ListForCourse(pageSize int, pageToken string, courseId string) ([]models.Certificate, string, error) {
	intermediateDb := InMemDb{}
	for _, cert := range rcv.records {
		if cert.CourseId == courseId {
			intermediateDb.records = append(intermediateDb.records, cert)
		}
	}
	if len(intermediateDb.records) == 0 {
		return nil, "", errors.New("No certificates were found for course")
	}
	return intermediateDb.List(pageSize, pageToken)
}

func (rcv *InMemDb) Delete(id string) {
	for k, cert := range rcv.records {
		if cert.Id == id {
			rcv.records[k], rcv.records[len(rcv.records)-1] = rcv.records[len(rcv.records)-1], rcv.records[k]
		}
	}
	rcv.records = rcv.records[0 : len(rcv.records)-2]
}

func (rcv *InMemDb) Connect(connectionString string) error {
	log.Println("initilazing local In-Memory Database...")
	rcv.init()
	log.Println("done!")
	return nil
}

func (rcv *InMemDb) Disconnect() {
	log.Println("flushing In-Memory Database...")
	rcv.records = nil
	log.Println("done!")
}

// func newIMD generates data for in-memory storage
func (rcv *InMemDb) init() {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Duration(2))
		rcv.records = append(rcv.records, models.Certificate{Id: fmt.Sprint(uuid.New()), UserId: fmt.Sprint(uuid.New()), CourseId: fmt.Sprint(uuid.New()), CreatedAt: time.Now()})
	}
}
