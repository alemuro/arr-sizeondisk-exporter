![GitHub Tag](https://img.shields.io/github/v/tag/alemuro/arr-sizeondisk-exporter)
![GitHub top language](https://img.shields.io/github/languages/top/alemuro/arr-sizeondisk-exporter)

# Arr SizeOnDisk exporter

Small and lightweight tool to export the size on disk of Radarr and Sonarr libraries.

## Usage

Environment variables:
- `RADARR_APIKEY`: The API key of the Radarr instance.
- `RADARR_HOST`: The URL of the Radarr instance.
- `SONARR_APIKEY`: The API key of the Sonarr instance.
- `SONARR_HOST`: The URL of the Sonarr instance.

### CLI

```bash
RADARR_APIKEY=apikey
RADARR_HOST=https://radarr...
SONARR_APIKEY=apikey
SONARR_HOST=https://sonarr...

$ go run cmd/cli/main.go
2024/07/14 11:33:42 Starting server on :9101
```

### Docker

```bash
$ docker run -e RADARR_APIKEY=apikey \
                -e RADARR_HOST=https://radarr... \
                -e SONARR_APIKEY=apikey \
                -e SONARR_HOST=https://sonarr... \
                -p 9101:9101 \
                ghcr.io/alemuro/arr-sizeondisk-exporter:latest

2024/07/14 11:33:42 Starting server on :9101
```

### Kubernetes

```hcl
module "exporter" {
  source  = "alemuro/expose-service-ingress/kubernetes"
  version = "1.23.0"

  namespace = "default"

  name           = "sizeondisk-exporter"
  image          = "ghcr.io/alemuro/arr-sizeondisk-exporter:latest"
  container_port = "9101"
  environment_variables = {
    RADARR_HOST   = "http://radarr"
    RADARR_APIKEY = var.radarr_api_key
    SONARR_HOST   = "http://sonarr"
    SONARR_APIKEY = var.sonarr_api_key
  }
  annotations = {
    service = {
      "k8s.grafana.com/scrape"                 = "true"
      "k8s.grafana.com/metrics.portName"       = "http"
      "k8s.grafana.com/metrics.scheme"         = "http"
      "k8s.grafana.com/metrics.path"           = "/metrics"
      "k8s.grafana.com/metrics.scrapeInterval" = "5m"
    }
  }
}
```

## Metrics

The service exposes a metric called `size_on_disk` which returns the size in bytes of a
movie/series season with the following labels:
- Category: `radarr` or `sonarr`.
- Collection: The name of the movie/series collection where the movie or season belongs to.
- Title: The title of the movie or series season.
