package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/IdeaEvolver/cutter-pkg/clog"
	"github.com/IdeaEvolver/cutter-pkg/service"
	"github.com/IdeaEvolver/cutter-status-dashboard/server"
	"github.com/IdeaEvolver/cutter-status-dashboard/status"
	"github.com/kelseyhightower/envconfig"

	"contrib.go.opencensus.io/integrations/ocsql"
)

type Config struct {
	DbHost     string `envconfig:"DB_HOST" required:"true"`
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

	statusStore := status.New(cfg.PlatformHealthcheck, cfg.FulfillmentHealthcheck, cfg.CrmHealthcheck, cfg.StudyHealthcheck, db)

	scfg := &service.Config{
		Addr:                fmt.Sprintf(":%s", cfg.PORT),
		ShutdownGracePeriod: time.Second * 10,
		MaxShutdownTime:     time.Second * 30,
	}

	handler := &server.Handler{
		Statuses: statusStore,
	}
	s := server.New(scfg, handler)

	clog.Infof("listening on %s", s.Addr)
	fmt.Println(s.ListenAndServe())
}
