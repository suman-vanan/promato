# 🍅 Promato

CLI for Prometheus, _inspired by_ [PromLens](https://github.com/prometheus/promlens)

## Usage

```shell
# Configure Promato to use a Prometheus HTTP API at a specified URL (for example, http://localhost:9090)
promato config --url http://localhost:9090

# Explore all metrics
promato explore

# Explore a specific named metric series and its labels, for example the series named "prometheus_http_requests_total"
promato explore prometheus_http_requests_total	
```