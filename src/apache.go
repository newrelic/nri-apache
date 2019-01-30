package main

import (
	"net/url"
	"time"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	StatusURL    string `default:"http://127.0.0.1/server-status?auto" help:"Apache status-server URL."`
	CABundleFile string `default:"" help:"Alternative Certificate Authority bundle file"`
	CABundleDir  string `default:"" help:"Alternative Certificate Authority bundle directory"`
}

const (
	integrationName    = "com.newrelic.apache"
	integrationVersion = "1.1.0"

	defaultHTTPTimeout = time.Second * 1
)

var args argumentList

func main() {
	log.Debug("Starting Apache integration")
	defer log.Debug("Apache integration exited")

	integration, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	fatalIfErr(err)
	u, err := url.Parse(args.StatusURL)
	fatalIfErr(err)
	entity, err := integration.Entity(integrationName, u.Hostname())
	fatalIfErr(err)

	if args.All() || args.Inventory {
		log.Debug("Fetching data for '%s' integration", integrationName+"-inventory")
		fatalIfErr(setInventory(entity.Inventory, u))
	}

	if args.All() || args.Metrics {
		log.Debug("Fetching data for '%s' integration", integrationName+"-metrics")
		ms := entity.NewMetricSet("ApacheSample")
		provider := &Status{
			CABundleDir:  args.CABundleDir,
			CABundleFile: args.CABundleFile,
			HTTPTimeout:  defaultHTTPTimeout,
		}
		fatalIfErr(getMetricsData(provider, ms))
	}

	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
