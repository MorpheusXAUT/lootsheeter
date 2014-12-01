// log
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/kdar/factorlog"
)

var (
	logger = factorlog.New(os.Stdout, factorlog.NewStdFormatter("[%{Date} %{Time}] {%{SEVERITY}:%{File}/%{Function}:%{Line}} %{SafeMessage}"))
)

func SetupLogger() {
	logger.SetMinMaxSeverity(factorlog.Severity(1<<uint(*debugLevelFlag)), factorlog.PANIC)
}

func WebLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		logger.Debugf("WebLogger: [%s] %q {%s} - %s ", r.Method, r.RequestURI, name, time.Since(start))
	})
}
