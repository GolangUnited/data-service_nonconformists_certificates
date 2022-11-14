package api

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// errGeneral is an error when we don't want to share details
	errGeneral = status.New(codes.Internal, "Something went wrong").Err()

	// errIncorrectInput is an error when we can't parse input(usually string to UUID)
	errIncorrectInput = status.New(codes.Internal, "Incorrect input").Err()

	// errCertNotFound is an error when no certificate was found in database
	errCertNotFound = status.New(codes.NotFound, "Certificate not found").Err()

	// errCertAlreadyExists is an error when certificate for user and course was found in database
	errCertAlreadyExists = status.New(codes.AlreadyExists, "Cert for this user for this course was already created").Err()
)
