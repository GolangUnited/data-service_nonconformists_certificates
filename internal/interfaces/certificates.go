// Package db contains DB interface
// and it's implementations for different backends.

package interfaces

import "golang-united-certificates/internal/models"

type CertificatesRepos interface {
	Connect(connectionString string) error
	GetById(id string) (result models.Certificate, err error)
	Create(cert *models.Certificate) error
	List(listOptions models.ListOptions) (result []models.Certificate, NextPageToken string, err error)
	Delete(id string) error
	Disconnect()
}
