package chinadns

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/proxy"

	"github.com/caddyserver/caddy"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

func init() {
	caddy.RegisterPlugin("chinadns", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {

	var dbfile string
	var upstream proxy.Upstream
	var err error

	for c.Next() {
		if !c.Dispenser.Args(&dbfile) {
			return c.ArgErr()
		}

		upstream, err = proxy.NewStaticUpstream(&c.Dispenser)
		if err != nil {
			return plugin.Error("chinadns", err)
		}
	}

	geoip_db, err := maxminddb.Open(dbfile)
	if err != nil {
		return plugin.Error("chinadns", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Chinadns{
			Next:  next,
			Proxy: proxy.Proxy{Upstreams: &[]proxy.Upstream{upstream}},
			Geoip: geoip_db,
		}
	})

	return nil
}
