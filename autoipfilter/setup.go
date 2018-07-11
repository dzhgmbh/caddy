package autoipfilter

import (
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/mholt/caddy/caddytls"
	"net/http"
	"net"
	"github.com/pkg/errors"
	"fmt"
)

var privateIpFilterMiddleWare func(next httpserver.Handler) httpserver.Handler

type AutoIpFilter struct {
	Next httpserver.Handler
}

func init() {
	caddy.RegisterPlugin("autoipfilter", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})

	privateIpFilterMiddleWare = func(next httpserver.Handler) httpserver.Handler {
		return AutoIpFilter{Next: next}
	}
}

func setup(c *caddy.Controller) error {
	cfg := httpserver.GetConfig(c)

	if !caddytls.HostQualifies(cfg.Host()) {
		cfg.AddMiddleware(privateIpFilterMiddleWare)
	}

	return nil
}

func (handler AutoIpFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	clientIp := net.ParseIP(host)

	if clientIp == nil {
		return http.StatusInternalServerError, errors.New(fmt.Sprintf("unable to parse client IP: %v", host))
	}

	for _, privateSubnet := range caddytls.PrivateIpBlocks {
		if privateSubnet.Contains(clientIp) {
			return handler.Next.ServeHTTP(w, r)
		}
	}

	return http.StatusForbidden, nil
}
