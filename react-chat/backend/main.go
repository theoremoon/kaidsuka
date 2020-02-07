package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/theoremoon/kaidsuka/react-chat/backend/repository"
	"github.com/theoremoon/kaidsuka/react-chat/backend/server"
	"github.com/theoremoon/kaidsuka/react-chat/backend/service"
)

//go:generate statik -src ./db -f

func run() error {
	repo, err := repository.New("database.db")
	if err != nil {
		return err
	}
	if err := repo.Setup(); err != nil {
		return err
	}

	service, err := service.New(repo)
	if err != nil {
		return err
	}
	server := server.New(service)

	return server.Start(":8000")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
