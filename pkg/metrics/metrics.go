package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Config Metrics 中间件配置
type Config struct {
	// SkipPaths 跳过统计的路径
	SkipPaths []string
}

// ==================== HTTP 指标 ====================

var (
	// HTTPRequestTotal HTTP 请求总数（Counter）
	HTTPRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestDuration HTTP 请求耗时分布（Histogram）
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency distributions",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// HTTPRequestsInFlight 当前处理中的请求数（Gauge）
	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)
)

// ==================== 数据库指标 ====================

var (
	// DBConnectionsActive 数据库活跃连接数
	DBConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	// DBConnectionsIdle 数据库空闲连接数
	DBConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	// DBQueryDuration 数据库查询耗时
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query latency distributions",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)
)

// ==================== Redis 指标 ====================

var (
	// RedisCommandsTotal Redis 命令总数
	RedisCommandsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_commands_total",
			Help: "Total number of Redis commands",
		},
		[]string{"command", "status"},
	)

	// RedisRequestDuration Redis 请求耗时
	RedisRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_request_duration_seconds",
			Help:    "Redis request latency distributions",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"command"},
	)
)

// ==================== 应用指标 ====================

var (
	// AppStartTime 应用启动时间
	AppStartTime = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_start_time_seconds",
			Help: "Application start time as unix timestamp",
		},
	)

	// AppInfo 应用信息
	AppInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_info",
			Help: "Application information",
		},
		[]string{"version", "go_version"},
	)
)