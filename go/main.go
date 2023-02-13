package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
	"log"
	"main.go/infra"
	"math/rand"
	"net/http"
	"time"
)

var tracer trace.Tracer

const cKeyMetrics = "custom_metrics"

type Metrics struct {
	customCnt *prometheus.Metric
	customDur *prometheus.Metric
}

// Needed by echo-contrib so echo can register and collect these metrics
func (m *Metrics) MetricList() []*prometheus.Metric {
	return []*prometheus.Metric{
		// ADD EVERY METRIC HERE!
		m.customCnt,
		m.customDur,
	}
}

// Creates and populates a new Metrics struct
// This is where all the prometheus metrics, names and labels are specified
func NewMetrics() *Metrics {
	return &Metrics{
		customCnt: &prometheus.Metric{
			Name:        "custom_total",
			Description: "Custom counter events.",
			Type:        "counter_vec",
			Args:        []string{"label_one", "label_two"},
		},
		customDur: &prometheus.Metric{
			Name:        "custom_duration_seconds",
			Description: "Custom duration observations.",
			Type:        "histogram_vec",
			Args:        []string{"label_one", "label_two"},
		},
	}
}

func (m *Metrics) AddCustomMetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(cKeyMetrics, m)
		return next(c)
	}
}

func (m *Metrics) IncCustomCnt(labelOne, labelTwo string) {
	labels := prom.Labels{"label_one": labelOne, "label_two": labelTwo}
	m.customCnt.MetricCollector.(*prom.CounterVec).With(labels).Inc()
}

func (m *Metrics) ObserveCustomDur(labelOne, labelTwo string, d time.Duration) {
	labels := prom.Labels{"label_one": labelOne, "label_two": labelTwo}
	m.customCnt.MetricCollector.(*prom.HistogramVec).With(labels).Observe(d.Seconds())
}


func main() {
	viper.SetDefault("JAEGER_ENDPOINT_EXPORTER", "http://localhost:14268/api/traces")
	viper.SetDefault("CLIENT_ENDPOINT", "http://localhost:8080/otel")
	viper.AutomaticEnv()

	ot := infra.NewOtel()
	ot.ServiceName = "goOtel"
	ot.ExporterEndpoint = viper.GetString("JAEGER_ENDPOINT_EXPORTER")
	tracer = ot.GetTracer()

	m := NewMetrics()

	e := echo.New()
	p := prometheus.NewPrometheus("echo", nil, m.MetricList())
	p.Use(e)
	e.Use(m.AddCustomMetricsMiddleware)

	e.Use(otelecho.Middleware(ot.ServiceName))

	e.GET("/test", func(c echo.Context) error {
		options := []string{"metric", "trace", "log"}
		rand.Seed(time.Now().UnixNano())
		label := options[rand.Intn(len(options))]
		metrics := c.Get(cKeyMetrics).(*Metrics)
		metrics.IncCustomCnt(label, "value")
		ctx := baggage.ContextWithoutBaggage(c.Request().Context())

		// 1 - simulating the execution of something
		ctx, executionSomething := tracer.Start(ctx, "execution-something")
		time.Sleep(time.Millisecond * 100)
		executionSomething.End()

		// 2 - Request Kotlin service
		ctx, httpCall := tracer.Start(ctx, "http-call")
		test := getTestOtel()
		httpCall.End()

		// 3 - Show result
		ctx, showResult := tracer.Start(ctx, "show-result")
		res := c.String(http.StatusOK, test)
		showResult.End()
		return res
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func getTestOtel() string {
	client := resty.New()
	resp, err := client.R().EnableTrace().Get(viper.GetString("CLIENT_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	return string(resp.Body())
}
