package api

import (
	"context"
	"golang-united-certificates/internal/interfaces"
	"golang-united-certificates/internal/models"
	"log"

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
	if !isUUID(req.GetId()) {
		return nil, errIncorrectInput
	}
	cert, err := srv.Database.GetById(req.GetId())
	if err != nil {
		return nil, errCertNotFound
	}
	return &GetResponse{Certificate: writeApiCert(cert)}, nil
}

func (srv *GRPCServer) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	if !isUUID(req.GetUserId()) || !isUUID(req.GetCourseId()) {
		return nil, errIncorrectInput
	}
	cert := models.Certificate{
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
	return &CreateResponse{Certificate: writeApiCert(cert)}, nil
}

func (srv *GRPCServer) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	if cid := req.GetCourseId(); !isUUID(cid) && len(cid) != 0 {
		return nil, errIncorrectInput
	}
	if cid := req.GetUserId(); !isUUID(cid) && len(cid) != 0 {
		return nil, errIncorrectInput
	}
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
		resp[k] = writeApiCert(cert)
	}
	return &ListResponse{Certificates: resp}, nil
}

func (srv *GRPCServer) Delete(ctx context.Context, req *DeleteRequest) (*emptypb.Empty, error) {
	if !isUUID(req.GetId()) {
		return nil, errIncorrectInput
	}
	cert := models.Certificate{
		ID: req.GetId(),
	}
	err := srv.Database.Delete(&cert)
	if err != nil {
		return &emptypb.Empty{}, errGeneral
	}
	return &emptypb.Empty{}, nil
}

func writeApiCert(cert models.Certificate) *Cert {
	return &Cert{
		Id:        cert.ID,
		UserId:    cert.UserId,
		CourseId:  cert.CourseId,
		CreatedAt: timestamppb.New(cert.CreatedAt),
		IsDeleted: cert.IsDeleted,
	}
}

// returns nil if UUID and
func isUUID(s string) bool {
	if _, err := uuid.Parse(s); err != nil {
		log.Printf("gon an incorrect input. want UUID, got %s", s)
		return false
	}
	return true
}
