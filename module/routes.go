package module

import (
	as "github.com/fuzailAhmad123/test_report/module/activity_snap" //activity_snap
	rsrvc "github.com/fuzailAhmad123/test_report/module/report"     //rsrvc
	"github.com/fuzailAhmad123/test_report/module/types"
	"github.com/go-chi/chi/v5"
)

func Route(hr *types.HTTPAPIResource) chi.Router {
	r := chi.NewRouter()

	// routes for activity report
	r.Get("/", rsrvc.GetActivityReportController(hr))
	r.Post("/aggregate-activity-snap", as.ActivitySnapAggregateController(hr))

	return r
}
