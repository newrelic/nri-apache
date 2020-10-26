package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/persist"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	StatusURL        string `default:"http://127.0.0.1/server-status?auto" help:"Apache status-server URL."`
	CABundleFile     string `default:"" help:"Alternative Certificate Authority bundle file"`
	CABundleDir      string `default:"" help:"Alternative Certificate Authority bundle directory"`
	RemoteMonitoring bool   `default:"false" help:"Identifies the monitored entity as 'remote'. In doubt: set to true."`
	ValidateCerts    bool   `default:"true" help:"If the status URL is HTTPS with a self-signed certificate, set this to false if you want to avoid certificate validation"`
	ShowVersion      bool   `default:"false" help:"Print build information and exit"`
}

const (
	integrationName = "com.newrelic.apache"

	defaultHTTPTimeout = time.Second * 1

	entityRemoteType = "server"
	httpProtocol     = `http`

	httpsProtocol    = `https`
	httpDefaultPort  = `80`
	httpsDefaultPort = `443`
)

var (
	args               argumentList
	integrationVersion = "0.0.0"
	gitCommit          = ""
	buildDate          = ""
)

func main() {
	log.Debug("Starting Apache integration")
	defer log.Debug("Apache integration exited")

	i, err := createIntegration()
	fatalIfErr(err)

	if args.ShowVersion {
		fmt.Printf(
			"New Relic %s integration Version: %s, Platform: %s, GoVersion: %s, GitCommit: %s, BuildDate: %s\n",
			strings.Title(strings.Replace(integrationName, "com.newrelic.", "", 1)),
			integrationVersion,
			fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			runtime.Version(),
			gitCommit,
			buildDate)
		os.Exit(0)
	}

	log.SetupLogging(args.Verbose)

	e, err := entity(i, args.StatusURL, args.RemoteMonitoring)
	fatalIfErr(err)

	if args.HasInventory() {
		log.Debug("Fetching data for '%s' integration", integrationName+"-inventory")
		fatalIfErr(setInventory(e.Inventory))
	}

	if args.HasMetrics() {
		log.Debug("Fetching data for '%s' integration", integrationName+"-metrics")

		hostname, port, err := parseStatusURL(args.StatusURL)
		fatalIfErr(err)

		ms := metricSet(e, "ApacheSample", hostname, port, args.RemoteMonitoring)
		provider := &Status{
			CABundleDir:   args.CABundleDir,
			CABundleFile:  args.CABundleFile,
			HTTPTimeout:   defaultHTTPTimeout,
			ValidateCerts: args.ValidateCerts,
		}
		fatalIfErr(getMetricsData(provider, ms))
	}

	fatalIfErr(i.Publish())
}

func entity(i *integration.Integration, statusURL string, remote bool) (*integration.Entity, error) {
	if remote {
		hostname, port, err := parseStatusURL(statusURL)
		if err != nil {
			return nil, err
		}
		n := fmt.Sprintf("%s:%s", hostname, port)
		return i.Entity(n, entityRemoteType)
	}

	return i.LocalEntity(), nil
}

func metricSet(e *integration.Entity, eventType, hostname, port string, remote bool) *metric.Set {
	if remote {
		return e.NewMetricSet(
			eventType,
			metric.Attr("hostname", hostname),
			metric.Attr("port", port),
		)
	}

	return e.NewMetricSet(
		eventType,
		metric.Attr("port", port),
	)
}

// parseStatusURL will extract the hostname and the port from the apache status URL.
func createIntegration() (*integration.Integration, error) {
	cachePath := os.Getenv("NRIA_CACHE_PATH")
	if cachePath == "" {
		return integration.New(integrationName, integrationVersion, integration.Args(&args))
	}

	l := log.NewStdErr(args.Verbose)
	s, err := persist.NewFileStore(cachePath, l, persist.DefaultTTL)
	if err != nil {
		return nil, err
	}

	return integration.New(integrationName, integrationVersion, integration.Args(&args), integration.Storer(s), integration.Logger(l))
}

// parseStatusURL will extract the hostname and the port from the nginx status URL.
func parseStatusURL(statusURL string) (hostname, port string, err error) {
	u, err := url.Parse(statusURL)
	if err != nil {
		return
	}

	if !isHTTP(u) {
		err = errors.New("unsupported protocol scheme")
		return
	}

	hostname = u.Hostname()

	if hostname == "" {
		err = errors.New("http: no Host in request URL")
		return
	}

	if u.Port() != "" {
		port = u.Port()
	} else if u.Scheme == httpsProtocol {
		port = httpsDefaultPort
	} else {
		port = httpDefaultPort
	}
	return
}

// isHTTP is checking if the URL is http/s protocol.
func isHTTP(u *url.URL) bool {
	return u.Scheme == httpProtocol || u.Scheme == httpsProtocol
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
