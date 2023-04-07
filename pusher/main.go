package main

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/orange-cloudfoundry/gomod_exporter/common"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
	"github.com/alecthomas/kingpin/v2"
)

var (
	pushGwURL       = kingpin.Flag("pushgw-url", "Push gateway url").Required().String()
	pushGwSkipSSL   = kingpin.Flag("pushgw-unsecure", "Skip SSL verify").Bool()
	metricJobName   = kingpin.Flag("metric-job-name", "name seen by prometheus as job_name").Default("gomod").String()
	metricNS        = kingpin.Flag("metric-namespace", "metric prefix namespace").Default("gomod").String()
	projectURL      = kingpin.Flag("project-url", "Git target project to analyze").Required().String()
	projectUsername = kingpin.Flag("project-user", "(optional) username for git authentication").String()
	projectPassword = kingpin.Flag("project-password", "(optional) password for git authentication").String()
	projectDir      = kingpin.Flag("project-dir", "(optional) use given directory instead of cloning project").String()
	fake            = kingpin.Flag("fake", "(optional) do not push metrics, only prints on stdout").Bool()
	logLevel        = kingpin.Flag("log-level", "Log level").Default("info").String()
	logJSON         = kingpin.Flag("log-json", "Log in JSON").Bool()
	logNoColor      = kingpin.Flag("log-no-color", "Disable log coloring").Bool()
)

func main() {
	kingpin.Version(version.Print("gomod-pusher"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	config := NewConfig()
	common.InitLogs(&config.BaseConfig)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: config.gwSkipSSL,
	}
	metrics := common.NewMetrics(config.metricNS)
	analyzer := common.NewAnalyzer(&config.BaseConfig, metrics)
	pusher := push.New(config.gwURL, config.metricJobName)

	project := config.BaseConfig.Projects[0]
	if err := analyzer.ProcessProject(&project); err != nil {
		log.Warnf("unable to analyze project: %s", err)
		log.Warnf("failure will be reported in pushed metrics")
	}

	if *fake {
		gathering, err := metrics.Registry.Gather()
		if err != nil {
			log.Errorf("unable to gather metrics: %s", err)
			os.Exit(1)
		}
		for _, cMF := range gathering {
			_, err = expfmt.MetricFamilyToText(os.Stdout, cMF)
			if err != nil {
				log.Errorf("unable to output metric: %s", err)
				os.Exit(1)
			}
		}
		os.Exit(0)
	}

	pusher.Gatherer(metrics.Registry)
	pusher.Grouping("project", project.URL)
	if err := pusher.Add(); err != nil {
		log.Errorf("unable to push data to gateway: %s", err)
		os.Exit(1)
	}
}

// Local Variables:
// ispell-local-dictionary: "american"
// End:
