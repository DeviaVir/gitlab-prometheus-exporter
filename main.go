package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

func getEnv(name string) string {
	envValue, ok := os.LookupEnv(name)
	if ok {
		return envValue
	}
	panic(fmt.Sprintf("Missing environment variable: %s", name))
}

func getEnvDefault(name string, defaultVal string) string {
	envValue, ok := os.LookupEnv(name)
	if ok {
		return envValue
	}
	return defaultVal
}

func setGauge(name string, help string, callback func() float64) {
	gaugeFunc := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "gitlab",
		Subsystem: "api",
		Name:      name,
		Help:      help,
	}, callback)
	prometheus.MustRegister(gaugeFunc)
}

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	gitlabAPIURL := getEnv("GITLAB_API")
	gitlabToken := getEnv("GITLAB_TOKEN")
	listendAddr := getEnvDefault("HTTP_LISTENADDR", ":9111")

	// create gitlab client
	git, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabAPIURL))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// grab license
	license, _, err := git.License.GetLicense()
	if err != nil {
		log.Fatalf("Failed to grab license")
	}

	// set gauges
	setGauge("license_active_users", "License active users", func() float64 {
		return float64(license.ActiveUsers)
	})
	setGauge("license_overage", "Users outside of the license", func() float64 {
		return float64(license.Overage)
	})
	setGauge("license_user_limit", "Number of active users allowed inside the license", func() float64 {
		return float64(license.UserLimit)
	})
	http.Handle("/metrics", promhttp.Handler())
	logrus.Info("Now listening on ", listendAddr)
	logrus.Fatal(http.ListenAndServe(listendAddr, nil))
}
