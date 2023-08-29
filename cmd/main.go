package main

import (
	"bwg2/config"
	"bwg2/internal/tsfl/handlers"
	"bwg2/internal/tsfl/repository"
	"bwg2/internal/tsfl/usecase"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
)

func main() {
	fmt.Println(uuid.New().String())
	viperConf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	conf, err := config.ParseConfig(viperConf)
	if err != nil {
		log.Fatal(err)
	}

	rep, err := repository.NewRepository(context.Background(), conf)
	if err != nil {
		log.Fatal(err)
	}
	service := usecase.NewService(rep)

	serv := handlers.NewFiberServer(conf, service)

	err = serv.Run(conf)
	if err != nil {
		log.Fatal(err)
	}

}
