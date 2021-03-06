package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/eldeal/observation-controller/config"
	"github.com/eldeal/observation-controller/mapper"
	"github.com/eldeal/observation-controller/observations"
	"github.com/gorilla/mux"
)

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}

// Observations handles observation requests
func Observations(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		datasetID := vars["dataset_id"]
		edition := vars["edition"]
		version := vars["version"]

		//validate dataset details and get metadata
		//TODO: move client init to service pkg
		dsCli := dataset.NewAPIClient(cfg.APIURL)
		v, err := dsCli.GetVersion(ctx, "", "", "", "", datasetID, edition, version)
		if err != nil {
			log.Error(ctx, "failed to get version details", err)
			setStatusCode(r, w, err)
			return
		}

		//request observations (URL params to start, via a form later on)
		obsCli := observations.NewAPIClient(cfg.APIURL)
		obsDoc, err := obsCli.Get(ctx, datasetID, edition, version, r.URL.Query())
		if err != nil {
			log.Error(ctx, "failed to get version details", err)
			setStatusCode(r, w, err)
			return
		}

		//format response
		obs := mapper.Observation{
			Name: r.URL.EscapedPath(),
		}
		model := mapper.WithVersion(ctx, obs, v, cfg)
		if err := model.ParseObservationDetails(ctx, obsDoc); err != nil {
			setStatusCode(r, w, err)
			return
		}

		b, err := json.Marshal(model)
		if err != nil {
			setStatusCode(r, w, err)
			return
		}

		if _, err = w.Write(b); err != nil {
			log.Error(ctx, "failed to write bytes for http response", err)
			setStatusCode(r, w, err)
			return
		}
		return
	}
}
