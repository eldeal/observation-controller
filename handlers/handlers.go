package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/eldeal/observation-controller/config"
	"github.com/eldeal/observation-controller/mapper"
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

// Bulletin handles bulletin requests
func Observations(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		//vars := mux.Vars(r)
		// datasetID := vars["dataset_id"]
		// edition := vars["edition"]
		// version := vars["version"]

		obs := mapper.Observation{Name: r.URL.EscapedPath()}
		model := mapper.Blank(ctx, obs, cfg)

		b, err := json.Marshal(model)
		if err != nil {
			setStatusCode(r, w, err)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			log.Error(ctx, "failed to write bytes for http response", err)
			setStatusCode(r, w, err)
			return
		}
		return
	}
}
