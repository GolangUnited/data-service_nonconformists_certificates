// Package db contains DB interface
// and it's implementations for different backends.

package interfaces

import "golang-united-certificates/internal/models"

type CertificatesRepos interface {
	// Connect handles connection to database by given connection ID
	Connect(connectionString string) error

	// GetById returns certificate object by given ID.
	// Returns an empty object and NotFound error if no object is found.
	GetById(id string) (result models.Certificate, err error)

	// Create appends given certificate to database.
	// Fills given object with data which is unnecessary to append it to db(eg. creation timestamp)
	Create(cert *models.Certificate) error

	// List filter up all records by given set of options and filters.
	// Supports pagination with limit and offset. List do not checks
	// options for it's validity. If no data found - returns empty array of certificates.
	List(listOptions models.ListOptions) ([]models.Certificate, error)

	// Delete removes certificate with given ID from database
	Delete(id string) error

	// Disconnect closes all connections to database
	Disconnect()
}
