package main

import (
	// "context"
	"context"
	"fmt"

	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/varun-singhal-omniful/oms-service/Init"
	"github.com/varun-singhal-omniful/oms-service/router"
	// Init "github.com/varun-singhal-omniful/oms-service/init"
)

func main() {
	// fmt.Println("edfws")
	server := http.InitializeServer(":8080", 0, 0, 70)
	context := context.Background()
	err := router.Initialize(context, server)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	Init.InitializeDB(context)
	err = server.StartServer("OMS")
	if err != nil {
		fmt.Println("Error in starting the server")
		return
	}
	fmt.Println("Server started")
}
