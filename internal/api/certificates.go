package api

import (
	"context"
	"errors"

	"golang-united-certificates/internal/interfaces"
	"golang-united-certificates/internal/models"

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
	cert, err := srv.Database.GetCertById(req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
	}
	return &GetResponse{Certificate: WriteApiCert(cert)}, err
}

func (srv *GRPCServer) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	if srv.Database.IsCertExistsByUserAndCourse(req.UserId, req.CourseId) {
		return &CreateResponse{}, status.New(codes.AlreadyExists, errors.New("Cert for this user for this course was already created").Error()).Err()
	}
	cert, err := srv.Database.Create(req.UserId, req.CourseId)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &CreateResponse{Certificate: WriteApiCert(cert)}, nil
}

func (srv *GRPCServer) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	data, npt, err := srv.Database.List(int(req.PageSize), req.PageToken)
	resp := make([]*Cert, len(data))
	for k, cert := range data {
		resp[k] = WriteApiCert(cert)
	}
	if err != nil {
		err = status.New(codes.Internal, err.Error()).Err()
	}
	return &ListResponse{Certificates: resp, NextPageToken: npt}, err
}

func (srv *GRPCServer) ListForUser(ctx context.Context, req *ListForUserRequest) (*ListResponse, error) {
	data, npt, err := srv.Database.ListForUser(int(req.PageSize), req.PageToken, req.UserId)
	resp := make([]*Cert, len(data))
	for k, cert := range data {
		resp[k] = WriteApiCert(cert)
	}
	if err != nil {
		err = status.New(codes.Internal, err.Error()).Err()
	}
	return &ListResponse{Certificates: resp, NextPageToken: npt}, err
}
func (srv *GRPCServer) ListForCourse(ctx context.Context, req *ListForCourseRequest) (*ListResponse, error) {
	data, npt, err := srv.Database.ListForCourse(int(req.PageSize), req.PageToken, req.CourseId)
	resp := make([]*Cert, len(data))
	for k, cert := range data {
		resp[k] = WriteApiCert(cert)
	}
	if err != nil {
		err = status.New(codes.Internal, err.Error()).Err()
	}
	return &ListResponse{Certificates: resp, NextPageToken: npt}, err
}
func (srv *GRPCServer) Delete(ctx context.Context, req *DeleteRequest) (*emptypb.Empty, error) {
	srv.Database.Delete(req.Id)
	return &emptypb.Empty{}, nil
}

func WriteApiCert(cert models.Certificate) *Cert {
	return &Cert{Id: cert.Id, UserId: cert.UserId, CourseId: cert.CourseId, CreatedAt: timestamppb.New(cert.CreatedAt)}
}
