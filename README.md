# IBM Storage Insights Prometheus Exporter
Prometheus exporter for IBM Storage Insights

## Overview
IBM Storage Insights exposes REST APIs over which performance and capacity metrics of different storage systems and their components can be fetched. This application is a *Prometheus exporter* whose aim is to bridge those metrics into Prometheus. It runs as a server process that translates Storage Insights metrics into a format that Prometheus can scrape. The exporter architecture is shown below.
![image](https://github.com/user-attachments/assets/7b398e89-6bdd-4fd2-9ce0-cee2b9ea12a1)
Fetching and translating the metrics from Storage Insights happens synchronously for every scrape from Prometheus. Each time Prometheus scrapes metrics, the exporter calls Storage Insights Metric API passing the system details and the authentication token. The authentication token is valid for 15 minutes, after which a fresh token is fetched by the exporter automatically before making API calls.

## Prerequisites
Software and experience prerequisites for using the Prometheus Exporter:
- ### Software requirements
  The exact steps have been tested on MacOS. However, the overall procedure should be portable to any modern Linux distribution or Windows.
  You will need to have:
  - A functioning Prometheus monitoring setup
  - A functioning Go language setup
- ### Required experience
  You will need a basic understanding of Prometheus-based monitoring.
  
  You will also need to understand the basics of software development in Go.

## Storage Insights REST API
The exporter calls the following Storage Insights REST APIs
Name | URL | Method | GET parameters | Description
---|---|---|---|---
Token API | /restapi/v1/{tenant_uuid}/token | POST | - | Creates an API token for tenant user
Metric API | /restapi/v1/tenants/{tenant_uuid}/storage-systems/metrics | GET | metric types | Returns capacity and performance metric values for all storage systems of a given tenant

### Authentication
Authentication is performed by the exporter by passing an authentication token in the header <code>x-api-token</code>.

## Steps to export IBM Storage Insights metrics of a system to Prometheus
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
