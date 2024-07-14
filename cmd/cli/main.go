package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alemuro/arr-sizeondisk-exporter/internal/consts"
	"github.com/alemuro/arr-sizeondisk-exporter/internal/exporters"
	"github.com/alemuro/arr-sizeondisk-exporter/internal/logger"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()

	exporter := exporters.NewExporter()
	exporter.AddRadarrProvider(os.Getenv(consts.EnvRadarrApiKey), os.Getenv(consts.EnvRadarrHost))
	exporter.AddSonarrProvider(os.Getenv(consts.EnvSonarrApiKey), os.Getenv(consts.EnvSonarrHost))

	r := prometheus.NewRegistry()
	r.MustRegister(exporter)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	logger.Info(fmt.Sprintf("Starting server on port %s", consts.DefaultPort))

	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", consts.DefaultPort), nil))
}
