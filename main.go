package main

import (
	"context"
	"crypto/tls"
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
)

type Config struct {
	REDIS_URL              string `envconfig:"REDIS_URL"`
	PlatformEndpoint       string `envconfig:"PLATFORM_ENDPOINT" required:"false"`
	FulfillmentHealthcheck string `envconfig:"FULFILLMENT_ENDPOINT" required:"false"`
	CrmHealthcheck         string `envconfig:"CRM_ENDPOINT" required:"false"`
	StudyHealthcheck       string `envconfig:"STUDY_ENDPOINT" required:"false"`

	PORT string `envconfig:"PORT"`
}

func main() {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		clog.Fatalf("config: %s", err)
	}

	pool := status.InitRedis(cfg.REDIS_URL)
	statusStore := status.New(pool)

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
		Platform:    cfg.PlatformEndpoint,
		Fulfillment: cfg.FulfillmentHealthcheck,
		Crm:         cfg.CrmHealthcheck,
		Study:       cfg.StudyHealthcheck,
	}

	handler := &server.Handler{
		Statuses:     statusStore,
		Healthchecks: healthchecksClient,
	}

	s := server.New(scfg, handler)

	ctx := context.Background()
	go handler.AllChecks(ctx)

	clog.Infof("listening on %s", s.Addr)
	fmt.Println(s.ListenAndServe())
}
