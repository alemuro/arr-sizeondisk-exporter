package exporters

import (
	"fmt"

	"github.com/alemuro/arr-sizeondisk-exporter/internal/logger"
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

type Exporter struct {
	sonarrcli  *sonarr.Sonarr
	radarrcli  *radarr.Radarr
	sizeOnDisk *prometheus.Desc
}

// NewExporter creates a new Exporter
func NewExporter() *Exporter {
	return &Exporter{
		sonarrcli:  nil,
		radarrcli:  nil,
		sizeOnDisk: prometheus.NewDesc("size_on_disk", "Size on disk for the item", []string{"collection", "title", "category"}, nil),
	}
}

// AddRadarrProvider adds a Radarr provider to the Exporter
func (e *Exporter) AddRadarrProvider(apiKey, host string) {
	e.radarrcli = radarr.New(starr.New(apiKey, host, 0))
	logger.Info(fmt.Sprintf("Radarr provider added to host %s", host))
}

// AddSonarrProvider adds a Sonarr provider to the Exporter
func (e *Exporter) AddSonarrProvider(apiKey, host string) {
	e.sonarrcli = sonarr.New(starr.New(apiKey, host, 0))
	logger.Info(fmt.Sprintf("Sonarr provider added to host %s", host))
}

// Describe sends the super-set of all possible descriptors of metrics
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.sizeOnDisk
}

// Collect is called by the Prometheus registry when collecting metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// RADARR
	if e.radarrcli != nil {
		movies, err := e.getMovieStatistics()
		if err != nil {
			logger.Error(fmt.Sprintf("Error getting Radarr statistics: %s", err))
			return
		}

		for _, movie := range movies {
			ch <- prometheus.MustNewConstMetric(e.sizeOnDisk, prometheus.GaugeValue, float64(movie.Size), movie.Title, movie.Title, "radarr")
		}
	}

	// SONARR
	if e.sonarrcli != nil {
		seasons, err := e.getSeriesStatistics()
		if err != nil {
			logger.Error(fmt.Sprintf("Error getting Sonarr statistics: %s", err))
			return
		}

		for _, season := range seasons {
			title := fmt.Sprintf("%s (%s)", season.Title, season.Season)
			ch <- prometheus.MustNewConstMetric(e.sizeOnDisk, prometheus.GaugeValue, float64(season.Size), season.Title, title, "sonarr")
		}
	}
}

// getMovieStatistics returns the size on disk for each movie
func (e *Exporter) getMovieStatistics() (mlist []RadarrMovie, err error) {
	movies, err := e.radarrcli.GetMovie(0)

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
		logger.Debug(fmt.Sprintf("[Radarr] %s - %d", movie.Title, movie.SizeOnDisk))
	}

	return mlist, nil
}

// getSeriesStatistics returns the size on disk for each season
func (e *Exporter) getSeriesStatistics() (seasons []SonarrSeason, err error) {
	series, err := e.sonarrcli.GetSeries(0)

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
			logger.Debug(fmt.Sprintf("[Sonarr] %s (Season %d) - %d", serie.Title, season.SeasonNumber, season.Statistics.SizeOnDisk))
		}
	}

	return seasons, nil
}
