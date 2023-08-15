package main

import (
	"log"

	"github.com/techieasif/go-echo-postgres-micro-service/internal/database"
	"github.com/techieasif/go-echo-postgres-micro-service/internal/server"
)

func main(){
	dbClient, err := database.NewDatabaseClient()
	if err != nil {
		log.Fatalf("Fail to init DB: %s", err)
	}
	srv := server.NewEchoServer(dbClient)
	if err := srv.Start(); err != nil{
		log.Fatal(err.Error())
	}

}