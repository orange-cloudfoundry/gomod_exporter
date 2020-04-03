package main

import (
	"github.com/gorilla/mux"
	"github.com/orange-cloudfoundry/gomod_exporter/common"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
)

var (
	configFile = kingpin.Flag("config", "Configuration file path").Required().File()
)

func main() {
	kingpin.Version(version.Print("gomod-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	config := NewConfig(*configFile)
	common.InitLogs(&config.BaseConfig)

	metrics := common.NewMetrics(config.Exporter.Namespace)
	analyzer := common.NewAnalyzer(&config.BaseConfig, metrics)
	analyzer.RunForever(config.Exporter.intervalDuration)

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	if (config.Web.SSLCertPath != "") && (config.Web.SSLKeyPath != "") {
		log.Infof("serving https on %s", config.Web.Listen)
		panic(http.ListenAndServeTLS(config.Web.Listen, config.Web.SSLCertPath, config.Web.SSLKeyPath, router))
	}
	log.Infof("serving http on %s", config.Web.Listen)
	panic(http.ListenAndServe(config.Web.Listen, router))
}

// Local Variables:
// ispell-local-dictionary: "american"
// End:
