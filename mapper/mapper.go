package mapper

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-observation-api/models"
	"github.com/eldeal/observation-controller/config"
)

type Observation struct {
	// TODO define
	Name string `json:"name"`
}

type ObservationModel struct {
	// TODO define
	Name              string `json:"model-name"`
	ReleaseDate       string `json:"release_date"`
	FixedDimensions   []Dim
	Observations      []ObsResponse
	TotalObservations int
	Unit              string
}

type Dim struct {
	Name   string
	Link   string
	Option string
}

type ObsResponse struct {
	Dim
	Label string
	Value string
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

func (m *ObservationModel) ParseObservationDetails(ctx context.Context, doc models.ObservationsDoc) error {
	m.TotalObservations = doc.TotalObservations
	m.Unit = doc.UnitOfMeasure

	for name, opt := range doc.Dimensions {
		m.FixedDimensions = append(m.FixedDimensions, Dim{
			Name:   name,
			Link:   opt.LinkObject.URL,
			Option: opt.LinkObject.ID,
		})
	}

	if len(doc.Observations) == 1 {
		m.Observations = append(m.Observations, ObsResponse{Value: doc.Observations[0].Observation})
		return nil
	}

	for _, o := range doc.Observations {
		if len(o.Dimensions) != 1 {
			return errors.New("observation values must not contain multiple dimension options")
		}

		for name, dim := range o.Dimensions {

			m.Observations = append(m.Observations, ObsResponse{
				Dim: Dim{
					Name: name,
					Link: dim.HRef,
				},
				Label: dim.Label,
				Value: o.Observation,
			})
		}

	}
	return nil
}
