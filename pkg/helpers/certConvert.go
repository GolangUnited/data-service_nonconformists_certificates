package helpers

import (
	"time"

	"github.com/IlyaKhalizov/golang-united-certificates/pkg/api"
	"github.com/IlyaKhalizov/golang-united-certificates/pkg/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func WriteApiCert(cert models.Certificate) *api.Cert {
	return &api.Cert{CertId: cert.CertId, UserId: cert.UserId, CourseId: cert.CourseId, CreatedAt: timestamppb.New(cert.CreatedAt)}
}
func GetApiCert(apiCert *api.Cert) models.Certificate {
	return models.Certificate{CertId: apiCert.CertId, UserId: apiCert.UserId, CourseId: apiCert.CourseId, CreatedAt: time.Unix(apiCert.CreatedAt.Seconds, int64(apiCert.CreatedAt.Nanos))}
}
