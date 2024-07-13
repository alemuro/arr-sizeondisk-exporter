package exporters

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"golift.io/starr"
	"golift.io/starr/radarr"
	"golift.io/starr/sonarr"
)

type RadarrMovie struct {
	Title string
	Size  int64
}

type SonarrSeason struct {
	Title  string
	Season string
	Size   int64
}

type ExporterConfigProvider struct {
	APIKey string
	Host   string
}

type ExporterConfig struct {
	Radarr ExporterConfigProvider
	Sonarr ExporterConfigProvider
}

type Exporter struct {
	scli       *sonarr.Sonarr
	rcli       *radarr.Radarr
	sizeOnDisk *prometheus.Desc
}

func NewExporter(config ExporterConfig) *Exporter {
	return &Exporter{
		scli:       sonarr.New(starr.New(config.Sonarr.APIKey, config.Sonarr.Host, 0)),
		rcli:       radarr.New(starr.New(config.Radarr.APIKey, config.Radarr.Host, 0)),
		sizeOnDisk: prometheus.NewDesc("size_on_disk", "Size on disk for the item", []string{"collection", "title", "category"}, nil),
	}
}
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.sizeOnDisk
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// RADARR
	movies, err := e.getMovieStatistics()
	if err != nil {
		log.Printf("Error getting statistics: %s", err)
		return
	}

	for _, movie := range movies {
		ch <- prometheus.MustNewConstMetric(e.sizeOnDisk, prometheus.GaugeValue, float64(movie.Size), movie.Title, movie.Title, "radarr")
	}

	// SONARR
	seasons, err := e.getSeriesStatistics()
	if err != nil {
		log.Printf("Error getting statistics: %s", err)
		return
	}

	for _, season := range seasons {
		title := fmt.Sprintf("%s (%s)", season.Title, season.Season)
		ch <- prometheus.MustNewConstMetric(e.sizeOnDisk, prometheus.GaugeValue, float64(season.Size), season.Title, title, "sonarr")
	}
}

func (e *Exporter) getMovieStatistics() (mlist []RadarrMovie, err error) {
	movies, err := e.rcli.GetMovie(0)

	if err != nil {
		return nil, err
	}

	for _, movie := range movies {
		if movie.SizeOnDisk == 0 {
			continue
		}
		mlist = append(mlist, RadarrMovie{
			Title: movie.Title,
			Size:  movie.SizeOnDisk,
		})
		log.Printf("[Radarr] %s - %d", movie.Title, movie.SizeOnDisk)
	}

	return mlist, nil
}

func (e *Exporter) getSeriesStatistics() (seasons []SonarrSeason, err error) {
	series, err := e.scli.GetSeries(0)

	if err != nil {
		return nil, err
	}

	for _, serie := range series {
		for _, season := range serie.Seasons {
			if season.Statistics.SizeOnDisk == 0 {
				continue
			}
			seasons = append(seasons, SonarrSeason{
				Title:  serie.Title,
				Season: fmt.Sprintf("Season %d", season.SeasonNumber),
				Size:   season.Statistics.SizeOnDisk,
			})
			log.Printf("[Sonarr] %s (Season %d) - %d", serie.Title, season.SeasonNumber, season.Statistics.SizeOnDisk)
		}
	}

	return seasons, nil
}
