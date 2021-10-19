package mapper

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/eldeal/observation-controller/config"
)

type Observation struct {
	// TODO define
	Name string `json:"name"`
}

type ObservationModel struct {
	// TODO define
	Name        string `json:"model-name"`
	ReleaseDate string `json:"release_date"`
}

func Blank(ctx context.Context, obs Observation, cfg config.Config) ObservationModel {
	var model ObservationModel
	model.Name = obs.Name
	return model
}

func WithVersion(ctx context.Context, obs Observation, v dataset.Version, cfg config.Config) ObservationModel {
	var model ObservationModel
	model.Name = obs.Name

	model.ReleaseDate = v.ReleaseDate
	return model
}
