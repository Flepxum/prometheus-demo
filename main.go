package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// 1. Counter：只增不减
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP requests total",
		},
		[]string{"method", "path"},
	)

	// 2. Gauge：可增可减
	memoryGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
	)
	memoryGauge2 = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_memory_usage_bytes_labeled",
			Help: "Memory usage",
		},
		[]string{"service", "instance"},
	)

	// 3. Histogram：统计数据分布
	requestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration distribution",
			Buckets: prometheus.DefBuckets,
		},
	)

	// 4. Summary：统计分位数
	requestLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "http_request_latency_seconds",
			Help: "HTTP request latency quantiles",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
	)
)

func init() {
	prometheus.MustRegister(
		requestCounter,
		memoryGauge,
		memoryGauge2,
		requestDuration,
		requestLatency,
	)
}

// 定时采集内存使用量
func collectMemory() {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			memoryGauge.Set(float64(m.Alloc))
			memoryGauge2.WithLabelValues("my_service", "instance_1").Set(float64(m.Alloc) / 1024 / 1024) // 转换为 MB
		}
	}()
}

// 业务 HTTP Handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// 模拟业务处理耗时
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

	// 计算请求耗时
	duration := time.Since(start).Seconds()

	// Counter：请求数 +1
	requestCounter.WithLabelValues(
		r.Method,
		r.URL.Path,
	).Inc()

	// Histogram：记录请求耗时
	requestDuration.Observe(duration)

	// Summary：记录请求耗时
	requestLatency.Observe(duration)

	fmt.Fprintf(w, "Hello Prometheus!\n")
}

func main() {
	collectMemory()

	// 业务接口
	http.HandleFunc("/hello", helloHandler)

	// Prometheus 指标接口
	http.Handle(
		"/metrics",
		promhttp.Handler(),
	)

	log.Println("Server listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
