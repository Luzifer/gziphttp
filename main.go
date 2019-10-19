package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	httphelper "github.com/Luzifer/go_helpers/v2/http"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		Listen         string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		ServeDir       string `flag:"serve-dir,d" default:"." description:"Directory to serve files from"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func init() {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("gziphttp %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(cfg.ServeDir)))

	var handler http.Handler = http.DefaultServeMux
	handler = httphelper.GzipHandler(handler)
	handler = httphelper.NewHTTPLogHandler(handler)

	log.WithError(http.ListenAndServe(cfg.Listen, handler)).Error("HTTP server ended")
}
