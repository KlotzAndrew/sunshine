package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsResponseTime prometheus.Summary
)

func main() {
	httpRequestsResponseTime = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "db",
		Name:      "response_time_seconds",
		Help:      "Request response times",
	})
	prometheus.MustRegister(httpRequestsResponseTime)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/basic", echo.WrapHandler(prometheus.InstrumentHandler("basic", instrumentedHandler())))
	e.GET("/dbwrite", echo.WrapHandler(prometheus.InstrumentHandler("dbWrite", dbHandler())))
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func instrumentedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Basic instrumentation!"))
	})
}

func dbHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		res := writeToDb()

		w.Write([]byte(res))
	})
}

func writeToDb() string {
	start := time.Now()
	defer httpRequestsResponseTime.Observe(float64(time.Since(start).Seconds()))

	return "some db request!"
}
