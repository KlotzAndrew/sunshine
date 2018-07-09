package main

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	FastQuery prometheus.Summary
	SlowQuery prometheus.Summary
)

func main() {
	dbInstrument(&FastQuery, "fast_query")
	dbInstrument(&SlowQuery, "slow_query")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/basic", echo.WrapHandler(prometheus.InstrumentHandler("basic", instrumentedHandler())))
	e.GET("/dbwrite", echo.WrapHandler(prometheus.InstrumentHandler("dbWrite", dbHandler())))
	e.GET("/dbwrite_slow", echo.WrapHandler(prometheus.InstrumentHandler("dbWriteSlow", dbSlowHandler())))
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Start server
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(port))
}

func dbInstrument(query *prometheus.Summary, name string) {
	hostname, err := os.Hostname()
	if err != nil {
		panic(1)
	}
	*query = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "db",
		Name:      "query_time_seconds",
		Help:      "Query response times",
		ConstLabels: prometheus.Labels{
			"query":    name,
			"service":  "main_api",
			"hostname": hostname,
		},
	})
	prometheus.MustRegister(*query)
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

func dbSlowHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		res := writeToDbSlow()

		w.Write([]byte(res))
	})
}

func writeToDb() string {
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	defer FastQuery.Observe(float64(time.Since(start).Seconds()))

	return "some db request!"
}

func writeToDbSlow() string {
	start := time.Now()
	time.Sleep(1000 * time.Millisecond)
	defer SlowQuery.Observe(float64(time.Since(start).Seconds()))

	return "some slow db request!"
}
