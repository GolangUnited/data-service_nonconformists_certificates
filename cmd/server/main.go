package main

import (
	"fmt"
	"log"
	"net"

	"golang-united-certificates/config"
	"golang-united-certificates/internal/api"
	"golang-united-certificates/internal/interfaces"
	db "golang-united-certificates/internal/repositories/certificates"

	"google.golang.org/grpc"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	var database interfaces.CertificatesRepos
	switch conf.DBType {
	case "inmem":
		database = new(db.InMemDb)
	case "psql":
		database = new(db.PgSql)
	default:
		log.Println("starting with default DB...")
		database = new(db.InMemDb)
	}
	database.Connect(conf.ConnectionString)
	defer database.Disconnect()
	srv := grpc.NewServer()
	grpcsrv := &api.GRPCServer{Database: database}
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
