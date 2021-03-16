package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/IdeaEvolver/cutter-pkg/client"
	"github.com/IdeaEvolver/cutter-pkg/clog"
	"github.com/IdeaEvolver/cutter-pkg/service"
	"github.com/IdeaEvolver/cutter-status-dashboard/healthchecks"
	"github.com/IdeaEvolver/cutter-status-dashboard/server"
	"github.com/IdeaEvolver/cutter-status-dashboard/status"
	"github.com/kelseyhightower/envconfig"
	"go.opencensus.io/plugin/ochttp"

	"contrib.go.opencensus.io/integrations/ocsql"
)

type Config struct {
	DbHost     string `envconfig:"DB_HOSTNAME" required:"true"`
	DbPort     string `envconfig:"DB_PORT" required:"true"`
	DbUsername string `envconfig:"DB_USERNAME" required:"true"`
	DbPassword string `envconfig:"DB_PASSWORD" required:"true"`
	DbName     string `envconfig:"DB_NAME" required:"true"`
	DbOpts     string `envconfig:"DB_OPTS" required:"false"`

	PlatformHealthcheck    string `envconfig:"PLATFORM_HEALTHCHECK" required:"false"`
	FulfillmentHealthcheck string `envconfig:"FULFILLMENT_HEALTHCHECK" required:"false"`
	CrmHealthcheck         string `envconfig:"CRM_HEALTHCHECK" required:"false"`
	StudyHealthcheck       string `envconfig:"STUDY_HEALTHCHECK" required:"false"`

	PORT string `envconfig:"PORT"`
}

func main() {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		clog.Fatalf("config: %s", err)
	}

	cs := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbUsername,
		cfg.DbName,
		cfg.DbPassword,
	)

	driverName, err := ocsql.Register(
		"postgres",
		ocsql.WithQuery(true),
		ocsql.WithQueryParams(true),
		ocsql.WithInstanceName("status-dashboard"),
	)
	if err != nil {
		clog.Fatalf("unable to register our ocsql driver: %v\n", err)
	}

	db, err := sql.Open(driverName, cs)
	if err != nil {
		clog.Fatalf("failed to connect to db")
	}

	clog.Infof("connected to postgres: %s:%s", cfg.DbHost, cfg.DbPort)

	statusStore := status.New(db)

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	internalClient := &http.Client{
		Transport: &ochttp.Transport{
			// Use Google Cloud propagation format.
			Propagation: &propagation.HTTPFormat{},
			Base:        customTransport,
		},
	}

	scfg := &service.Config{
		Addr:                fmt.Sprintf(":%s", cfg.PORT),
		ShutdownGracePeriod: time.Second * 10,
		MaxShutdownTime:     time.Second * 30,
	}

	healthchecksClient := &healthchecks.Client{
		Client:      client.New(internalClient),
		Platform:    cfg.PlatformHealthcheck,
		Fulfillment: cfg.FulfillmentHealthcheck,
		Crm:         cfg.CrmHealthcheck,
		Study:       cfg.StudyHealthcheck,
	}

	handler := &server.Handler{
		Healthchecks: healthchecksClient,
		Statuses:     statusStore,
	}
	s := server.New(scfg, handler)

	//ctx := context.Background()
	//go handler.AllChecks(ctx)

	clog.Infof("listening on %s", s.Addr)
	fmt.Println(s.ListenAndServe())
}
