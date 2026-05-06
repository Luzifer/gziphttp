// gzip capable very simple HTTP server
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	httphelper "github.com/Luzifer/go_helpers/http"
	"github.com/Luzifer/rconfig/v2"
	"github.com/sirupsen/logrus"
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

func initApp() (err error) {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return fmt.Errorf("parsing cli options: %w", err)
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("parsing log-level: %w", err)
	}
	logrus.SetLevel(l)

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	if cfg.VersionAndExit {
		fmt.Printf("gziphttp %s\n", version) //nolint:forbidigo // fine for version print
		os.Exit(0)
	}

	http.Handle("/", http.FileServer(http.Dir(cfg.ServeDir)))

	var handler http.Handler = http.DefaultServeMux
	handler = httphelper.GzipHandler(handler)
	handler = httphelper.NewHTTPLogHandler(handler)

	server := &http.Server{
		Addr:              cfg.Listen,
		Handler:           handler,
		ReadHeaderTimeout: time.Second,
	}

	logrus.WithField("addr", cfg.Listen).Info("gziphttp started")
	if err = server.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("running HTTP server")
	}
}
