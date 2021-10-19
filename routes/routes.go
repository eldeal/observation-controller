package routes

import (
	"context"
	"net/http"

	"github.com/eldeal/observation-controller/config"
	"github.com/eldeal/observation-controller/handlers"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Clients - struct containing all the clients for the controller
type Clients struct {
	HealthCheckHandler func(w http.ResponseWriter, req *http.Request)
}

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, c Clients) {
	log.Info(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(c.HealthCheckHandler)
	r.StrictSlash(true).Path("/datasets/{dataset_id}/editions/{edition}/versions/{version}/observations").Methods("GET").HandlerFunc(handlers.Observations(*cfg))
}
