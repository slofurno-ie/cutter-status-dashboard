module github.com/IdeaEvolver/cutter-status-dashboard

go 1.14

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	contrib.go.opencensus.io/integrations/ocsql v0.1.7
	github.com/IdeaEvolver/cutter-pkg/client v1.7.0
	github.com/IdeaEvolver/cutter-pkg/clog v1.1.6
	github.com/IdeaEvolver/cutter-pkg/cuterr v1.4.0
	github.com/IdeaEvolver/cutter-pkg/service v0.7.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/rs/cors v1.7.0
	go.opencensus.io v0.23.0
)
