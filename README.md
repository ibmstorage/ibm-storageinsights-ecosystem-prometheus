# IBM Storage Insights Prometheus Exporter
Prometheus exporter for IBM Storage Insights

# Prerequisites
Software and experience prerequisites for using the Prometheus Exporter:
- ## Software requirements
  The exact steps have been tested on MacOS. However, the overall procedure should be portable to any modern Linux distribution or Windows.
  You will need to have:
  - A functioning Prometheus monitoring setup
  - A functioning Go language setup
- ## Required experience
  You will need a basic understanding of Prometheus-based monitoring.
  
  You will also need to understand the basics of software development in Go.

# Steps to export IBM Storage Insights metrics of a system to Prometheus
- Prepare `config.json`. Provide the following details:
  - Your IBM ID
  - Your REST API Key
  - Your Storage Insights tenant UUID
  - The system UUID whose metrics you want to fetch
  - The list of metrics that you're interested in
- Start the exporter
```
    go run main.go --listen-address :8089 --config ../config.json
```
- Add exporter as target and set scrape interval in `prometheus.yml`
```
    static_configs:
      - targets: ["localhost:8089"]
```
- Start Prometheus
- Open Prometheus url in a browser
