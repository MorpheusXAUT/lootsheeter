// log
package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kdar/factorlog"
)

var (
	logger = factorlog.New(os.Stdout, factorlog.NewStdFormatter("[%{Date} %{Time}] {%{SEVERITY}:%{File}/%{Function}:%{Line}} %{SafeMessage}"))
)

func SetupLogger() {
	logger.SetMinMaxSeverity(factorlog.Severity(1<<uint(config.DebugLevel)), factorlog.PANIC)
}

func WebLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if config.DebugTemplates {
			templates = template.Must(template.New("").Funcs(TemplateFunctions(nil)).ParseGlob("web/template/*"))
		}

		inner.ServeHTTP(w, r)

		remoteAddr := r.Header.Get("X-Forwarded-For")

		if len(remoteAddr) > 0 {
			remoteAddrs := strings.Split(remoteAddr, ", ")
			if len(remoteAddrs) > 1 {
				r.RemoteAddr = remoteAddrs[0]
			} else {
				r.RemoteAddr = remoteAddr
			}
		}

		logger.Debugf("WebLogger: [%s] %s %q {%s} - %s ", r.Method, r.RemoteAddr, r.RequestURI, name, time.Since(start))
	})
}
