# Go + Prometheus + Grafana 监控栈

本目录提供 Docker Compose 配置。Go 应用镜像使用已经发布到 GHCR 的：
`ghcr.io/flepxum/prometheus-demo:latest`

Prometheus 和 Grafana 使用官方镜像，分别运行在独立容器中。

## 启动

```bash
cp .env.example .env
# 编辑 .env，把 GRAFANA_ADMIN_PASSWORD 改成自己的密码
docker compose pull
docker compose up -d
```

访问：

- Go 指标： http://localhost:8080/metrics
- Prometheus： http://localhost:9090
- Grafana： http://localhost:3000 （用户名 admin，密码为 .env 中设置的值）

## 第一次查看图表

先产生一些业务请求，否则 Histogram / Summary 的 count 还是 0：

```bash
for i in $(seq 1 30); do
  curl -s http://localhost:8080/hello > /dev/null
done
```

Grafana 会自动配置 Prometheus 数据源，并加载 `Go Prometheus 监控面板`。

## PromQL 示例

```promql
# 每秒请求数
sum(rate(http_requests_total[1m]))

# 请求总数
sum(http_requests_total)

# Histogram P99
histogram_quantile(0.99, sum by (le) (rate(http_request_duration_seconds_bucket[5m])))

# Go 堆内存 MiB
go_memstats_heap_alloc_bytes / 1024 / 1024

# 进程 RSS MiB
process_resident_memory_bytes / 1024 / 1024

# Goroutine 数量
go_goroutines

# 进程 CPU（单核百分比）
rate(process_cpu_seconds_total[1m]) * 100
```

## 常用命令

```bash
docker compose ps
docker compose logs -f
docker compose logs -f prometheus
docker compose logs -f grafana
docker compose down
```

`docker compose down` 不会删除持久化数据卷。若要连数据一起删除，才使用 `docker compose down -v`。

## 注意

- Prometheus 在 Compose 网络内使用 `go-app:8080` 抓取 Go 指标，不要写 `localhost:8080`。
- Grafana 在 Compose 网络内通过 `http://prometheus:9090` 访问 Prometheus。
- GHCR 镜像必须对当前机器可拉取；私有镜像需先 `docker login ghcr.io`。
- 如果你实际的 GHCR 仓库所有者不是 `flepxum`，请相应修改 `docker-compose.yml` 中的镜像地址。
