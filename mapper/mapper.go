package mapper

import (
	"context"

	"github.com/eldeal/observation-controller/config"
)

type Observation struct {
	// TODO define
	Name string `json:"name"`
}

type ObservationModel struct {
	// TODO define
	Name string `json:"model-name"`
}

func Blank(ctx context.Context, obs Observation, cfg config.Config) ObservationModel {
	var model ObservationModel
	model.Name = obs.Name
	return model
}
