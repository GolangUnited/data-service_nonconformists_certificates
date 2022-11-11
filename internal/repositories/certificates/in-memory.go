package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang-united-certificates/internal/models"

	"github.com/google/uuid"
)

// InMemDb is implementations of simple In-Memory storage
// based on simple array.
// USE IT ONLY FOR DEVELOPMENT PURPOSES
type InMemDb struct {
	records []models.Certificate
}

// GetById returns whole certificate with given ID
// If there is no certificate with given ID, it returns
// empty struct and NotFound error
func (rcv *InMemDb) GetById(id string) (models.Certificate, error) {
	for _, cert := range rcv.records {
		if cert.Id == id {
			return cert, nil
		}
	}
	return models.Certificate{}, errors.New("No cert was found")
}

// Create adds given certificate to database and
// fills up it's properties.
// Always returns nil as error
func (rcv *InMemDb) Create(cert *models.Certificate) error {
	cert.CreatedAt = time.Now()
	rcv.records = append(rcv.records, *cert)
	return nil
}

// List returns an array of certificates, filtered by given listOptions filter
// If there is no records found or offset is too high - an error is returned
func (rcv *InMemDb) List(listOptions models.ListOptions) ([]models.Certificate, error) {
	var fResult []models.Certificate
	for _, cert := range rcv.records {
		if filterIfDeleted(cert, listOptions.ShowDeleted) {
			if filterByUserID(cert, listOptions.UserId) {
				if filterByCourseID(cert, listOptions.CourseId) {
					fResult = append(fResult, cert)
				}
			}
		}
	}
	if len(fResult) == 0 {
		log.Println("no records found")
		return nil, errors.New("no records found")
	}
	if listOptions.Offset >= len(fResult) {
		return nil, errors.New("Incorrect page token")
	}

	var result []models.Certificate
	if listOptions.Offset+listOptions.Limit >= len(fResult) {
		result = append(result, fResult[listOptions.Offset:]...)
	} else {
		result = append(result, fResult[listOptions.Offset:listOptions.Offset+listOptions.Limit]...)
	}
	return result, nil
}

// filterByUserID returns true if certificate's UserId matches
// given one, or if no UserId passed
func filterByUserID(cert models.Certificate, uid string) bool {
	return cert.UserId == uid || uid == ""
}

// filterByCourseID returns true if certificate's CourseId matches
// given one, or if no CourseId passed
func filterByCourseID(cert models.Certificate, cid string) bool {
	return cert.CourseId == cid || cid == ""
}

// filterIfDeleted returns true if DeletedAt is zero(not filled)
func filterIfDeleted(cert models.Certificate, showDeleted bool) bool {
	return cert.DeletedAt.IsZero() || showDeleted
}

// Delete marks certificate in DB as deleted
// imitating soft delete from real DB
// by changing DeletedAt
// and DeletedBy properties.
// Always returns nil as error.
func (rcv *InMemDb) Delete(inCert *models.Certificate) error {
	for k, crt := range rcv.records {
		if crt.Id == inCert.Id && crt.DeletedAt.IsZero() {
			rcv.records[k].DeletedAt = time.Now()
			rcv.records[k].DeletedBy = inCert.DeletedBy
			break
		}
	}
	return nil
}

// Connect do dummy stuff and init DB
// Always returns nil as error
func (rcv *InMemDb) Connect(connectionString string) error {
	log.Println("initializing local In-Memory Database...")
	rcv.init()
	log.Println("done!")
	return nil
}

// Disconnect do dummy stuff and init DB
// Always returns nil as error
func (rcv *InMemDb) Disconnect() {
	log.Println("flushing In-Memory Database...")
	rcv.records = nil
	log.Println("done!")
}

// init generates data for in-memory storage
func (rcv *InMemDb) init() {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Duration(2))
		rcv.records = append(
			rcv.records,
			models.Certificate{
				Id:        fmt.Sprint(uuid.New()),
				UserId:    fmt.Sprint(uuid.New()),
				CourseId:  fmt.Sprint(uuid.New()),
				CreatedAt: time.Now(),
				CreatedBy: fmt.Sprint(uuid.New()),
			})
	}
}
