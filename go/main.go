package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
	"log"
	"main.go/infra"
	"net/http"
	"time"
)

var tracer trace.Tracer

func main() {
	ot := infra.NewOtel()
	ot.ServiceName = "goOtel"
	ot.ExporterEndpoint = "http://localhost:14268/api/traces"
	tracer = ot.GetTracer()

	e := echo.New()
	e.Use(otelecho.Middleware(ot.ServiceName))

	e.GET("/test", func(c echo.Context) error {
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
	resp, err := client.R().EnableTrace().Get("http://localhost:8080/otel")
	if err != nil {
		log.Fatal(err)
	}
	return string(resp.Body())
}
