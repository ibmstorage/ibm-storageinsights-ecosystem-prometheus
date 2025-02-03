# IBM Storage Insights Prometheus Exporter
Prometheus exporter for IBM Storage Insights

## Overview
Prometheus is the de-facto monitoring system for metrics that is well-suited for monitoring dynamic systems. It features a dimensional data model as well as a powerful query language and integrates instrumentation, gathering of metrics, service discovery, and alerting, all in one ecosystem.

This application is intended for storage system administrators who want to onboard metrics gathered by IBM Storage Insights onto Prometheus, so that they can leverage alerting and powerful query language, PromQL, offered by Prometheus.

IBM Storage Insights exposes REST APIs over which performance and capacity metrics of different storage systems and their components can be fetched. This application is a **Prometheus exporter** whose aim is to bridge those metrics into Prometheus. It runs as a server process that translates Storage Insights metrics into a format that Prometheus can scrape. The exporter architecture is shown below.
![image](https://github.com/user-attachments/assets/7b398e89-6bdd-4fd2-9ce0-cee2b9ea12a1)
Fetching and translating the metrics from Storage Insights happens synchronously for every scrape from Prometheus. Each time Prometheus scrapes metrics, the exporter calls Storage Insights Metric API passing the request details and the authentication token. The authentication token is valid for 15 minutes, after which a fresh token is fetched by the exporter automatically before making API calls.

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
Token API | /restapi/v1/tenants/{tenant_uuid}/token | POST | - | Creates an API token for tenant user
Metric API | /restapi/v1/tenants/{tenant_uuid}/storage-systems/metrics | GET | metric types | Returns capacity and performance metric values for all storage systems of a given tenant

### Authentication
Authentication is performed by the exporter by passing an authentication token in the header <code>x-api-token</code>. The authentication token is retrived by calling the Token API passing the API key in the header <code>x-api-key</code>. A tenant admin can generate API key for tenant users by logging into [Storage Insights GUI](https://www.ibm.com/docs/en/storage-insights?topic=configuring-user-access-management#rest_api__section_dyd_t1t_jzb__title__1).

## Steps to export IBM Storage Insights metrics of a system to Prometheus
- Clone the repository
  ```bash
  git clone https://github.com/ibmstorage/ibm-storageinsights-ecosystem-prometheus.git
  cd ibm-storageinsights-ecosystem-prometheus
  ```
  
- Prepare `config.json`.
<img width="361" alt="config" src="https://github.com/user-attachments/assets/dc25d8cd-de67-4889-bc7d-81a8f5814598" />

  Provide the following details:
  - Your REST API Key
  - Your Storage Insights tenant UUID
  - The list of metrics that you're interested in

- Build the code. The code is made up of three Go modules, build each of them.
  - Build `config` module
    ```bash
    cd config
    go mod tidy
    go build
    ```
  - Build `simetrics` module
    ```bash
    cd simetrics
    go mod tidy
    go build
    ```
  - Build `main` module
    ```bash
    cd main
    go mod tidy
    go build
    ```
  
- Start the exporter. From the root directory 
```
    go run ./main/main.go --listen-address :8089 --config config.json
```
- Add exporter as target and set scrape interval in `prometheus.yml`
```
    static_configs:
      - targets: ["localhost:8089"]
```
- Start Prometheus. Change directory to where you've downloaded and extracted Prometheus distribution.
  ```bash
  ./prometheus
  ```
- Open Prometheus url in a browser
  <img width="1482" alt="metrics" src="https://github.com/user-attachments/assets/762b702a-893e-4637-81d3-6ca1851a015b" />

## Extending Prometheus Exporter
The application is also intended to act as a guide to those developers who want to extend the functionality of the exporter. The primary mode of getting metrics from Storage Insights is via its [REST API](https://insights.ibm.com/restapi/docs/). There are a multiple of Metrics API that provide metrics from different systems and their components. The details of how to invoke Metrics API is illustrated in the exporter file `si-metrics.go` in the simetrics Go module. To fetch more metrics, it is advisable to follow these steps:

- Declare the metric types in the config module
- Create a new Go module to fetch the additional metrics
  ```bash
  go init work <module-name>
  ```
  Use the code in module `simetrics` as an illustration of how to use REST API to fetch metrics data from Storage Insights.
- Refer to the Storage Insights REST API docs to get the details of the API.
- Register the new metrics with Prometheus
- Push the metrics to Prometheus on every scrape, optionally adding tags for description.
- Build the modules

### Limitations
* IBM Storage Insights collects metrics every 5 mins from the storage systems. Therefore, it is suggested to scrape at a frequency not less than 1 minute. In case you hit the rate limit, it is advisable to increase the scraping interval.
* Metrics API allows fetching upto 3 metric types for each of the storage system in one REST API call. If more than 3 metric types are desired, it is advisable to make multiple REST API calls as appropriate. 

### License
This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
