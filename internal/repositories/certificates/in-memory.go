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

func (rcv *InMemDb) GetById(id string) (models.Certificate, error) {
	for _, cert := range rcv.records {
		if cert.Id == id {
			return cert, nil
		}
	}
	return models.Certificate{}, errors.New("No cert was found")
}

func (rcv *InMemDb) IsExistsForUserAndCourse(userId, courseId string) bool {
	for _, cert := range rcv.records {
		if cert.UserId == userId && cert.CourseId == courseId {
			return true
		}
	}
	return false
}

func (rcv *InMemDb) Create(cert *models.Certificate) error {
	cert.CreatedAt = time.Now()
	rcv.records = append(rcv.records, *cert)
	return nil
}

func (rcv *InMemDb) List(listOptions models.ListOptions) ([]models.Certificate, string, error) {
	var npt string

	pt, err := strconv.Atoi(listOptions.PageToken)
	if err != nil {
		return nil, npt, err
	}
	var fresult []models.Certificate
	for _, cert := range rcv.records {
		if filterByUserID(cert, listOptions.UserId) {
			if filterByCourseID(cert, listOptions.CourseId) {
				fresult = append(fresult, cert)
			}
		}
	}
	if len(fresult) == 0 {
		return nil, npt, errors.New("no records found")
	}
	if pt >= len(fresult) {
		return nil, npt, errors.New("Incorrect page token")
	}

	var result []models.Certificate
	if pt+listOptions.PageSize >= len(fresult) {
		result = append(result, fresult[pt:len(fresult)]...)
	} else {
		result = append(result, fresult[pt:pt+listOptions.PageSize]...)
		npt = strconv.Itoa(pt + listOptions.PageSize)
	}
	return result, npt, nil
}

func filterByUserID(cert models.Certificate, uid string) bool {
	return cert.UserId == uid || uid == ""
}
func filterByCourseID(cert models.Certificate, cid string) bool {
	return cert.CourseId == cid || cid == ""
}

func (rcv *InMemDb) Delete(id string) error {
	for k, cert := range rcv.records {
		if cert.Id == id {
			rcv.records[k], rcv.records[len(rcv.records)-1] = rcv.records[len(rcv.records)-1], rcv.records[k]
		}
	}
	rcv.records = rcv.records[0 : len(rcv.records)-2]
	return nil
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
