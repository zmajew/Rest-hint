package cont

import (
	"Rest-hint/config"
	"Rest-hint/n4"
	"Rest-hint/server"
	"Rest-hint/service"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(config.NewConfig)
	container.Provide(config.ConnectDatabase)
	container.Provide(n4.NewMovieRep)
	container.Provide(service.NewMovieService)
	container.Provide(server.NewServer)

	return container
}
