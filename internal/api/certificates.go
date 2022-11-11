package api

import (
	"context"
	"golang-united-certificates/internal/interfaces"
	"golang-united-certificates/internal/models"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	UnimplementedCertificatesServer
	Database interfaces.CertificatesRepos
}

func NewCertApi(db interfaces.CertificatesRepos) *GRPCServer {
	return &GRPCServer{
		Database: db,
	}
}

func (srv *GRPCServer) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	cert, err := srv.Database.GetById(req.Id)
	if err != nil {
		return nil, errCertNotFound
	}
	return &GetResponse{Certificate: WriteApiCert(cert)}, nil
}

func (srv *GRPCServer) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	cert := models.Certificate{
		Id:       uuid.New().String(),
		UserId:   req.GetUserId(),
		CourseId: req.GetCourseId(),
	}
	err := srv.Database.Create(&cert)
	if err != nil {
		switch err.Error() {
		case "AlreadyExists":
			return nil, errCertAlreadyExists
		default:
			return nil, errGeneral
		}
	}
	return &CreateResponse{Certificate: WriteApiCert(cert)}, nil
}

func (srv *GRPCServer) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	listOptions := models.ListOptions{
		Limit:       int(req.GetLimit()),
		Offset:      int(req.GetOffset()),
		UserId:      req.GetUserId(),
		CourseId:    req.GetCourseId(),
		ShowDeleted: req.GetShowDeleted(),
	}
	listOptions.SetDefaults()
	data, err := srv.Database.List(listOptions)
	if err != nil {
		return nil, errGeneral
	}
	resp := make([]*Cert, len(data))
	for k, cert := range data {
		resp[k] = WriteApiCert(cert)
	}
	return &ListResponse{Certificates: resp}, nil
}

func (srv *GRPCServer) Delete(ctx context.Context, req *DeleteRequest) (*emptypb.Empty, error) {
	cert := models.Certificate{
		Id:        req.GetId(),
		DeletedBy: req.GetDeletedBy(),
	}
	err := srv.Database.Delete(&cert)
	if err != nil {
		return &emptypb.Empty{}, errGeneral
	}
	return &emptypb.Empty{}, nil
}

func WriteApiCert(cert models.Certificate) *Cert {
	return &Cert{
		Id:        cert.Id,
		UserId:    cert.UserId,
		CourseId:  cert.CourseId,
		CreatedAt: timestamppb.New(cert.CreatedAt),
		CreatedBy: cert.CreatedBy,
		DeletedAt: timestamppb.New(cert.DeletedAt),
		DeletedBy: cert.DeletedBy,
	}
}
