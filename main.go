package main

import (
	"Rest-hint/cont"
	"Rest-hint/server"
	"log"
)

func main() {
	container := cont.BuildContainer()

	err := container.Invoke(func(server *server.Server) {
		server.Run()
	})

	if err != nil {
		log.Fatalln(err)
	}
}
