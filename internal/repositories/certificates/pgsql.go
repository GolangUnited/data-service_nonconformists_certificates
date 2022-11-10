package db

import (
	"errors"
	"golang-united-certificates/internal/models"
	"log"
	"strconv"

	"github.com/google/uuid"
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

func (rcv *PgSql) GetCertById(id string) (models.Certificate, error) {
	var cert models.Certificate
	err := rcv.db.Model(&model).Take(&cert, "id = ?", id).Error
	return cert, err
}

func (rcv *PgSql) IsCertExistsByUserAndCourse(userId, courseId string) bool {
	var ctRec int64
	rcv.db.Model(&model).Where("user_id = ? AND course_id = ?", userId, courseId).Count(&ctRec)
	return ctRec != 0
}

func (rcv *PgSql) Create(userId, courseId string) (models.Certificate, error) {
	cert := models.Certificate{Id: uuid.New().String(), UserId: userId, CourseId: courseId}
	err := rcv.db.Create(&cert).Error
	if err != nil {
		return models.Certificate{}, err
	}
	return cert, nil
}

func (rcv *PgSql) List(pageSize int, pageToken string) ([]models.Certificate, string, error) {
	var certs []models.Certificate
	var npt string
	mod := rcv.db.Model(&model)
	if pageSize == 0 {
		err := mod.Find(&certs).Error
		log.Println(certs)
		return certs, npt, err
	}
	ptInt, err := convertPageToken(pageToken)
	if err != nil {
		return certs, npt, err
	}
	var ctRec int64
	mod.Count(&ctRec)
	if ctRec < int64(ptInt) {
		log.Println("pageToken is incorrect")
		return certs, npt, errors.New("Incorrect page token")
	}
	err = mod.Offset(ptInt).Limit(pageSize).Find(&certs).Error
	if err != nil {
		return certs, npt, err
	}
	if ctRec > (int64(ptInt) + int64(pageSize)) {
		npt = strconv.Itoa(ptInt + pageSize)
	}
	return certs, npt, nil
}

func (rcv *PgSql) ListForUser(pageSize int, pageToken string, userId string) ([]models.Certificate, string, error) {
	var filtered PgSql
	filtered.db = rcv.db.Model(&model).Where("user_id = ?", userId)
	return filtered.List(pageSize, pageToken)
}

func (rcv *PgSql) ListForCourse(pageSize int, pageToken string, courseId string) ([]models.Certificate, string, error) {
	var filtered PgSql
	filtered.db = rcv.db.Model(&model).Where("course_id = ?", courseId)
	return filtered.List(pageSize, pageToken)
}

func (rcv *PgSql) Delete(id string) {
	err := rcv.db.Model(&model).Where("id = ?", id).Delete(models.Certificate{}).Error
	if err != nil {
		log.Println("can't delete")
	}
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
