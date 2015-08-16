package domains

import (
	"net/http"
	"strings"

	"github.com/codeblanche/golibs/logr"
	"github.com/labstack/echo"
)

type (
	hosts map[string]*echo.Echo
	names map[*http.Request]string
)

// Host sets the router for the given host names
func (h hosts) Host(e *echo.Echo, names ...string) {
	for _, name := range names {
		h[name] = e
	}
}

// ServeHTTP implementation of http.Handler
func (h hosts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")[0]
	parts := strings.Split(host, ".")
	for len(parts) > 0 {
		name := strings.Join(parts, ".")
		logr.Debug(name)
		if e, ok := h[name]; ok && e != nil {
			e.ServeHTTP(w, r)
			return
		}
		parts = parts[0 : len(parts)-1]
	}

	http.Error(w, "Domain not recognised", http.StatusNotFound)
}
