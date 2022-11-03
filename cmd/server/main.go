package main

import (
	"fmt"
	"log"
	"net"

	"github.com/koldyba/gu-certmgr/pkg/api"
	"github.com/koldyba/gu-certmgr/pkg/certificates"
	"github.com/koldyba/gu-certmgr/pkg/config"
	"github.com/koldyba/gu-certmgr/pkg/db"

	"google.golang.org/grpc"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	var database db.DB
	switch conf.DBType {
	case "inmem":
		database = new(db.InMemDb)
	default:
		database = new(db.InMemDb)
	}
	database.Connect(conf.ConnectionString)

	srv := grpc.NewServer()
	grpcsrv := &certificates.GRPCServer{Database: database}
	api.RegisterCertificatesServer(srv, grpcsrv)

	log.Printf("starting server on port %d", conf.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.Serve(listener); err != nil {
		log.Fatal(err)
	}
	log.Printf("server is up and ready!")
}
