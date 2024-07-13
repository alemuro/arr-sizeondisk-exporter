package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alemuro/arr-sizeondisk-exporter/internal/exporters"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()

	radarrCfg := exporters.ExporterConfigProvider{
		APIKey: os.Getenv("RADARR_APIKEY"),
		Host:   os.Getenv("RADARR_HOST"),
	}
	sonarrCfg := exporters.ExporterConfigProvider{
		APIKey: os.Getenv("SONARR_APIKEY"),
		Host:   os.Getenv("SONARR_HOST"),
	}
	cfg := exporters.ExporterConfig{
		Radarr: radarrCfg,
		Sonarr: sonarrCfg,
	}
	prometheus.MustRegister(exporters.NewExporter(cfg))

	log.Println("Starting server on :9101")

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
