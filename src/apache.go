package main

import (
	"errors"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/persist"
	"net/url"
	"os"
	"time"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	StatusURL        string `default:"http://127.0.0.1/server-status?auto" help:"Apache status-server URL."`
	CABundleFile     string `default:"" help:"Alternative Certificate Authority bundle file"`
	CABundleDir      string `default:"" help:"Alternative Certificate Authority bundle directory"`
	RemoteMonitoring bool   `default:"true" help:"Identifies the monitored entity as 'remote'. In doubt: set to true."`
}

const (
	integrationName    = "com.newrelic.apache"
	integrationVersion = "1.2.0"

	defaultHTTPTimeout = time.Second * 1

	entityRemoteType = "server"
	httpProtocol     = `http`

	httpsProtocol    = `https`
	httpDefaultPort  = `80`
	httpsDefaultPort = `443`
)

var args argumentList

func main() {
	log.Debug("Starting Apache integration")
	defer log.Debug("Apache integration exited")

	var i *integration.Integration
	var err error
	cachePath := os.Getenv("NRIA_CACHE_PATH")

	if cachePath == "" {
		i, err = integration.New(integrationName, integrationVersion, integration.Args(&args))
	} else {
		var storer persist.Storer

		logger := log.NewStdErr(args.Verbose)
		storer, err = persist.NewFileStore(cachePath, logger, persist.DefaultTTL)
		fatalIfErr(err)

		i, err = integration.New(integrationName, integrationVersion, integration.Args(&args),
			integration.Storer(storer), integration.Logger(logger))
	}

	fatalIfErr(err)
	log.SetupLogging(args.Verbose)

	hostname, port, err := parseStatusURL(args.StatusURL)
	fatalIfErr(err)

	e, err := entity(i, hostname, port)
	fatalIfErr(err)

	if args.HasInventory() {
		log.Debug("Fetching data for '%s' integration", integrationName+"-inventory")
		fatalIfErr(setInventory(e.Inventory))
	}

	if args.HasMetrics() {
		log.Debug("Fetching data for '%s' integration", integrationName+"-metrics")

		hostnameAttr := metric.Attr("hostname", hostname)
		portAttr := metric.Attr("port", port)

		ms := e.NewMetricSet("ApacheSample", hostnameAttr, portAttr)
		provider := &Status{
			CABundleDir:  args.CABundleDir,
			CABundleFile: args.CABundleFile,
			HTTPTimeout:  defaultHTTPTimeout,
		}
		fatalIfErr(getMetricsData(provider, ms))
	}

	fatalIfErr(i.Publish())
}

func entity(i *integration.Integration, hostname, port string) (*integration.Entity, error) {
	if args.RemoteMonitoring {
		n := fmt.Sprintf("%s:%s", hostname, port)
		return i.Entity(n, entityRemoteType)
	}

	return i.LocalEntity(), nil
}

// parseStatusURL will extract the hostname and the port from the apache status URL.
func parseStatusURL(statusURL string) (hostname, port string, err error) {
	u, err := url.Parse(statusURL)
	if err != nil {
		return
	}

	isHTTP := u.Scheme == httpProtocol || u.Scheme == httpsProtocol
	if !isHTTP {
		err = errors.New("unsupported protocol scheme")
		return
	}

	if u.Hostname() == "" {
		err = errors.New("http: no Host in request URL")
		return
	}

	hostname = u.Hostname()

	if u.Port() != "" {
		port = u.Port()
	} else if u.Scheme == httpsProtocol {
		port = httpsDefaultPort
	} else {
		port = httpDefaultPort
	}
	return
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
