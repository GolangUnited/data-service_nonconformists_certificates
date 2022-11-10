package api

import (
	"context"
	"errors"
	"golang-united-certificates/internal/interfaces"
	"golang-united-certificates/internal/models"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	UnimplementedCertificatesServer
	Database interfaces.CertificatesRepos
}

func (srv *GRPCServer) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	cert, err := srv.Database.GetById(req.Id)
	if err != nil {
		return nil, status.New(codes.NotFound, err.Error()).Err()
	}
	return &GetResponse{Certificate: WriteApiCert(cert)}, err
}

func (srv *GRPCServer) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	cert := models.Certificate{
		Id:       uuid.New().String(),
		UserId:   req.GetUserId(),
		CourseId: req.GetCourseId(),
	}
	listOptions := models.ListOptions{
		UserId:   cert.UserId,
		CourseId: cert.CourseId,
	}
	listOptions.SetDefaults()
	c, _ := srv.Database.List(listOptions)
	if len(c) != 0 {
		return &CreateResponse{}, status.New(codes.AlreadyExists, errors.New("Cert for this user for this course was already created").Error()).Err()
	}
	err := srv.Database.Create(&cert)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &CreateResponse{Certificate: WriteApiCert(cert)}, nil
}

func (srv *GRPCServer) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	listOptions := models.ListOptions{
		Limit:    int(req.GetLimit()),
		Offset:   int(req.GetOffset()),
		UserId:   req.GetUserId(),
		CourseId: req.GetCourseId(),
	}
	listOptions.SetDefaults()
	data, err := srv.Database.List(listOptions)
	resp := make([]*Cert, len(data))
	for k, cert := range data {
		resp[k] = WriteApiCert(cert)
	}
	if err != nil {
		err = status.New(codes.Internal, err.Error()).Err()
	}
	return &ListResponse{Certificates: resp}, err
}

func (srv *GRPCServer) Delete(ctx context.Context, req *DeleteRequest) (*emptypb.Empty, error) {
	err := srv.Database.Delete(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.New(codes.Internal, err.Error()).Err()
	}
	return &emptypb.Empty{}, nil
}

func WriteApiCert(cert models.Certificate) *Cert {
	return &Cert{Id: cert.Id, UserId: cert.UserId, CourseId: cert.CourseId, CreatedAt: timestamppb.New(cert.CreatedAt)}
}
