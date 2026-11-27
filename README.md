# Go Prometheus 全量监控栈

基于你提供的 `/metrics` 输出制作，包含 Go 自定义指标、Go Runtime/GC/Heap、Process 和 promhttp 指标的 Grafana 面板。

## 启动

```bash
cp .env.example .env
# 修改 .env 中的 Grafana 管理员密码
docker compose pull
docker compose up -d
```

访问：
- Go 指标：http://localhost:8080/metrics
- Prometheus：http://localhost:9090
- Grafana：http://localhost:3000 （用户名 admin，密码为 .env 中设置的值）

Grafana 会自动配置 Prometheus 数据源，并从 provisioning 自动加载完整 Dashboard。

## 产生业务样本

```bash
for i in $(seq 1 100); do
  curl -s http://localhost:8080/hello > /dev/null
done
```

## 常用 PromQL

```promql
# HTTP QPS
sum(rate(http_requests_total[1m]))

# Histogram P99
histogram_quantile(0.99, sum by (le) (rate(http_request_duration_seconds_bucket[5m])))

# Summary P99
http_request_latency_seconds{quantile="0.99"}

# Go Heap Alloc
go_memstats_heap_alloc_bytes

# Process RSS
process_resident_memory_bytes

# CPU（单核百分比）
rate(process_cpu_seconds_total[1m]) * 100

# Goroutine
go_goroutines
```

## 常用操作

```bash
docker compose ps
docker compose logs -f
docker compose down
```

`docker compose down` 默认保留 Prometheus 和 Grafana 数据卷；`docker compose down -v` 会删除数据卷。

## 说明

- Go 镜像使用 `ghcr.io/flepxum/prometheus-demo:latest`。
- Prometheus 通过 Compose 服务名 `go-app:8080` 抓取指标。
- Grafana 通过 `http://prometheus:9090` 访问 Prometheus。
- `app_memory_usage_bytes_labeled` 当前代码记录的是 `Alloc / TotalAlloc` 比值，不是内存占用百分比；面板按原始值展示，并在描述中注明。
- `process_cpu_seconds_total` 的 rate 乘 100 是单核百分比口径，多核时可能超过 100%。
