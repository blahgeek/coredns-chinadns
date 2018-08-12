package chinadns

import (
	"net"

	"github.com/miekg/dns"
	"golang.org/x/net/context"

	maxminddb "github.com/oschwald/maxminddb-golang"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/proxy"
	"github.com/coredns/coredns/request"
)

type Chinadns struct {
	Next  plugin.Handler
	Proxy proxy.Proxy
	Geoip *maxminddb.Reader
}

var logger = log.NewWithPlugin("chinadns")

func (cn Chinadns) isInsideChina(res *dns.Msg) bool {
	if len(res.Answer) == 0 {
		return false
	}

	for _, ans := range res.Answer {
		var ip net.IP

		if ans.Header().Rrtype == dns.TypeA {
			ip = ans.(*dns.A).A
		} else {
			continue
		}

		var record struct {
			Country struct {
				ISOCode string `maxminddb:"iso_code"`
			} `maxminddb:"country"`
		}
		err := cn.Geoip.Lookup(ip, &record)
		if err != nil {
			logger.Warningf("Unable to lookup IP: %v", err)
			return false
		}
		logger.Debugf("Lookup result: %v in %v", ip.String(), record.Country.ISOCode)

		if record.Country.ISOCode != "CN" {
			return false
		}
	}

	return true
}

func (cn Chinadns) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	state := request.Request{Req: r, Context: ctx, W: w}

	if state.QType() == dns.TypeA || state.QType() == dns.TypeAAAA {
		// lookup A record for both A and AAAA
		res, err := cn.Proxy.Lookup(state, state.Name(), dns.TypeA)

		if err == nil && res != nil && cn.isInsideChina(res) {
			m := new(dns.Msg)
			m.SetReply(r)
			m.Authoritative = res.Authoritative
			m.RecursionAvailable = m.RecursionAvailable
			// For AAAA requests, return empty response
			if state.QType() == dns.TypeA {
				m.Answer = res.Answer
			}
			w.WriteMsg(m)
			return dns.RcodeSuccess, nil
		}
	}

	return plugin.NextOrFailure(cn.Name(), cn.Next, ctx, w, r)
}

func (cn Chinadns) Name() string {
	return "chinadns"
}
