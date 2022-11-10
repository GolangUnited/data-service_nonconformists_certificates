package db

import (
	"errors"
	"golang-united-certificates/internal/models"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PgSql struct {
	db *gorm.DB
}

var model models.Certificate

func (rcv *PgSql) Connect(connectionString string) error {
	log.Println("starting using psql database")
	var err error
	rcv.db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	log.Println("connected successfully!")
	rcv.applyMigrations()

	return nil
}
func (rcv *PgSql) Disconnect() {
	db, err := rcv.db.DB()
	// TODO: add err return and correct handling
	if err != nil {
		db.Close()
	}
	db.Close()
}

func (rcv *PgSql) GetById(id string) (models.Certificate, error) {
	var cert models.Certificate
	err := rcv.db.Model(&model).Take(&cert, "id = ?", id).Error
	return cert, err
}

func (rcv *PgSql) IsExistsForUserAndCourse(userId, courseId string) bool {
	var ctRec int64
	rcv.db.Model(&model).Where("user_id = ? AND course_id = ?", userId, courseId).Count(&ctRec)
	return ctRec != 0
}

func (rcv *PgSql) Create(cert *models.Certificate) error {
	return rcv.db.Create(&cert).Error

}

func (rcv *PgSql) List(listOptions models.ListOptions) ([]models.Certificate, string, error) {
	var certs []models.Certificate
	var npt string
	mod := rcv.db.Model(&model)

	if listOptions.UserId != "" {
		mod = mod.Where("user_id = ?", listOptions.UserId)
	}
	if listOptions.CourseId != "" {
		mod = mod.Where("course_id = ?", listOptions.CourseId)
	}

	ptInt, err := convertPageToken(listOptions.PageToken)
	if err != nil {
		return certs, npt, err
	}

	var ctRec int64
	mod.Count(&ctRec)
	if ctRec < int64(ptInt) {
		log.Println("pageToken is incorrect")
		return certs, npt, errors.New("Incorrect page token")
	}
	err = mod.Offset(ptInt).Limit(listOptions.PageSize).Find(&certs).Error
	if err != nil {
		return certs, npt, err
	}
	if ctRec > (int64(ptInt) + int64(listOptions.PageSize)) {
		npt = strconv.Itoa(ptInt + listOptions.PageSize)
	}
	return certs, npt, nil
}

func (rcv *PgSql) Delete(id string) error {
	err := rcv.db.Model(&model).Where("id = ?", id).Delete(models.Certificate{}).Error
	if err != nil {
		log.Println("can't delete")
	}
	return err
}

func (rcv *PgSql) applyMigrations() {
	log.Println("applying automigrate(letting gorm do all the job)...")
	rcv.db.AutoMigrate(model)
}

func convertPageToken(pageToken string) (int, error) {
	if pageToken == "" {
		pageToken = "0"
	}
	ptInt, err := strconv.Atoi(pageToken)
	if err != nil {
		log.Println("unable to parse pageToken")
		err = errors.New("unable to parse pageToken")
	}
	if ptInt < 0 {
		log.Println("negative pageToken is not supported")
		err = errors.New("negative pageToken is not supported")
		ptInt = 0
	}
	return ptInt, err
}
