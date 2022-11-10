package db

import (
	"golang-united-certificates/internal/models"
	"log"

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

func (rcv *PgSql) List(listOptions models.ListOptions) ([]models.Certificate, error) {
	var certs []models.Certificate
	mod := rcv.db.Model(&model)

	if listOptions.UserId != "" {
		mod = mod.Where("user_id = ?", listOptions.UserId)
	}
	if listOptions.CourseId != "" {
		mod = mod.Where("course_id = ?", listOptions.CourseId)
	}
	err := mod.Offset(listOptions.Offset).Limit(listOptions.Limit).Find(&certs).Error
	if err != nil {
		return certs, err
	}
	return certs, nil
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
