package helpers

import (
	"time"

	"golang-united-certificates/pkg/api"
	"golang-united-certificates/pkg/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func WriteApiCert(cert models.Certificate) *api.Cert {
	return &api.Cert{Id: cert.Id, UserId: cert.UserId, CourseId: cert.CourseId, CreatedAt: timestamppb.New(cert.CreatedAt)}
}
func GetApiCert(apiCert *api.Cert) models.Certificate {
	return models.Certificate{Id: apiCert.Id, UserId: apiCert.UserId, CourseId: apiCert.CourseId, CreatedAt: time.Unix(apiCert.CreatedAt.Seconds, int64(apiCert.CreatedAt.Nanos))}
}
