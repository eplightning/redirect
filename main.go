package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type optionalString struct {
	set   bool
	value string
}

type appConfig struct {
	listen string
	host   optionalString
	path   optionalString
	scheme optionalString
	query  optionalString
	code   int
}

type handler struct {
	conf *appConfig
}

func main() {
	conf := loadConfig()
	srv := configureServer(conf)
	shutdownOnSignal(srv)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("listen error: %v", err)
	}
}

func configureServer(conf *appConfig) *http.Server {
	return &http.Server{
		Addr:    conf.listen,
		Handler: handler{conf: conf},
	}
}

func shutdownOnSignal(srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("shutdown error, forcefully closing: %v", err)
		}
	}()
}

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	location := url.URL{
		Host:     req.Host,
		Scheme:   "http",
		Path:     req.URL.Path,
		RawQuery: req.URL.RawQuery,
	}

	if h.conf.host.set {
		location.Host = h.conf.host.value
	}
	if h.conf.scheme.set {
		location.Scheme = h.conf.scheme.value
	}
	if h.conf.path.set {
		location.Path = h.conf.path.value
	}
	if h.conf.query.set {
		location.RawQuery = h.conf.query.value
	}

	w.Header().Set("Location", location.String())
	w.WriteHeader(h.conf.code)
}

func loadConfig() *appConfig {
	envString := func(key string, defaultVal string) string {
		if val, found := os.LookupEnv(key); found {
			return val
		}

		return defaultVal
	}

	envOptionalString := func(key string, defaultVal optionalString) optionalString {
		if val, found := os.LookupEnv(key); found {
			return optionalString{set: true, value: val}
		}

		return defaultVal
	}

	envInt := func(key string, defaultVal int) int {
		if str, found := os.LookupEnv(key); found {
			if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
				return val
			} else {
				log.Fatalf("appConfig error [%v]: %v not an int", key, str)
			}
		}

		return defaultVal
	}

	return &appConfig{
		listen: envString("LISTEN_ADDRESS", ":8080"),
		host:   envOptionalString("HOST_OVERRIDE", optionalString{}),
		path:   envOptionalString("PATH_OVERRIDE", optionalString{}),
		query:  envOptionalString("QUERY_OVERRIDE", optionalString{}),
		scheme: envOptionalString("SCHEME_OVERRIDE", optionalString{set: true, value: "https"}),
		code:   envInt("STATUS_CODE", 301),
	}
}
