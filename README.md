# si-prometheus-exporter
Prometheus exporter for Storage Insights

# Steps to export SI metrics of a system to Prometheus
- Prepare `config.json`. Provide following details:
  - Your IBM ID
  - Your REST API Key
  - Your tenant UUID
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
![image](https://media.github.ibm.com/user/392205/files/0ab631a6-a10d-4aae-b3c1-78bfa8e6d8f9)
