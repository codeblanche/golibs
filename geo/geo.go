package geo

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/alecthomas/geoip"
)

var (
	// TTL duration to cache ip country result
	TTL time.Duration

	cache map[*http.Request]*cacheItem
	geo   *geoip.GeoIP
	err   error
	mutex *sync.RWMutex
)

type (
	cacheItem struct {
		country *geoip.Country
		expires time.Time
	}
)

func init() {
	cache = make(map[*http.Request]*cacheItem, 0)
	geo, err = geoip.New()
	mutex = &sync.RWMutex{}
}

func middlewareHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if err != nil || hasCached(r) {
		return
	}

	c := geo.Lookup(net.ParseIP(r.RemoteAddr))
	e := time.Now().Add(TTL)

	mutex.Lock()
	defer mutex.Unlock()
	cache[r] = &cacheItem{
		country: c,
		expires: e,
	}
	return
}

func hasCached(r *http.Request) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	c, ok := cache[r]
	return ok && c.expires.After(time.Now())
}

// Middleware returns a middleware handlerfunc
func Middleware() http.HandlerFunc {
	return middlewareHandlerFunc
}

// Get the country detected for a request
func Get(r *http.Request) geoip.Country {
	mutex.RLock()
	defer mutex.RUnlock()
	if c, ok := cache[r]; ok && c.country != nil {
		return *c.country
	}
	return geoip.Country{}
}
