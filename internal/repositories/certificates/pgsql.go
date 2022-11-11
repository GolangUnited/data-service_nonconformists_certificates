package db

import (
	"errors"
	"golang-united-certificates/internal/models"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PgSql is implementation of CertificatesRepos based on
// PostgreSQL via gorm
type PgSql struct {
	db *gorm.DB
}

var model models.Certificate

// Connect connects to PSQL with given connection string
// with auto migrations.
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

// Disconnect closes connection to PostgeSQL DB
func (rcv *PgSql) Disconnect() {
	db, err := rcv.db.DB()
	// TODO: add err return and correct handling
	if err != nil {
		db.Close()
	}
	db.Close()
}

// GetById returns whole certificate with given ID
// If there is no certificate with given ID, it returns
// empty struct and NotFound error
func (rcv *PgSql) GetById(id string) (models.Certificate, error) {
	var cert models.Certificate
	res := rcv.db.Model(&model).Take(&cert, "id = ?", id)
	err := res.Error
	var ctRec int64
	if res.Count(&ctRec); ctRec == 0 {
		return models.Certificate{}, errors.New("No cert was found")
	}
	return cert, err
}

// Create adds given certificate to database and
// fills up it's properties. Returns "AlreadyExists" error if
// record for such user and course already in place.
func (rcv *PgSql) Create(cert *models.Certificate) error {
	listOptions := models.ListOptions{
		UserId:   cert.UserId,
		CourseId: cert.CourseId,
	}
	listOptions.SetDefaults()
	c, _ := rcv.List(listOptions)
	if len(c) != 0 {
		return errors.New("AlreadyExists")
	}
	return rcv.db.Create(&cert).Error
}

// List returns an array of certificates, filtered by given listOptions filter
// Returns empty array if no records was found
func (rcv *PgSql) List(listOptions models.ListOptions) ([]models.Certificate, error) {
	var certs []models.Certificate
	mod := rcv.db.Model(&model)

	if !listOptions.ShowDeleted {
		mod = mod.Where("deleted_at IS NULL")
	}
	if listOptions.UserId != "" {
		mod = mod.Where("user_id = ?", listOptions.UserId)
	}
	if listOptions.CourseId != "" {
		mod = mod.Where("course_id = ?", listOptions.CourseId)
	}
	err := mod.Offset(listOptions.Offset).Limit(listOptions.Limit).Find(&certs).Error
	if err != nil {
		log.Printf("an error occured while getting records from db: %s", err.Error())
		return certs, err
	}
	return certs, nil
}

// Delete removes certificate with given ID from database
// Always returns nil error
func (rcv *PgSql) Delete(cert *models.Certificate) error {
	toDel := rcv.db.Model(&model).Where("id = ? AND deleted_at IS NULL", cert.Id)
	cert.DeletedAt = time.Now()
	toDel.Updates(cert)
	return nil
}

// applyMigrations uses gorm AutoMigrate to create DB schema
func (rcv *PgSql) applyMigrations() {
	log.Println("applying migrations(letting gorm do all the job)...")
	rcv.db.AutoMigrate(model)
}
