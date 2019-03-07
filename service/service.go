package service

import (
	"Rest-hint/config"
	"Rest-hint/n4"
)

type MovieService struct {
	config *config.Config
	rep    *n4.MovieRep
}

func (service *MovieService) GiveHint(letters map[string]string) ([]n4.Movie, error) {
	if service.config.Enabled {
		return service.rep.GiveHint(letters)
	}
	return []n4.Movie{}, nil
}

func NewMovieService(config *config.Config, rep *n4.MovieRep) *MovieService {
	return &MovieService{config: config, rep: rep}
}
