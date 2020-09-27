> **Unmaintained. Checkout https://github.com/tuna/freedns-go instead**

---


ChinaDNS Plugin for CoreDNS
====

A CoreDns plugin that response to A-type queries for chinese websites only, like [ChinaDNS](https://github.com/shadowsocks/ChinaDNS), with several highlights:

- Simple, Hackable, Flexible
- Use GeoIP database to determine IP country
- Handle AAAA-type queries correctly, given that chinese users usually access IPv6 network via a foriegn proxy/tunnel
    - Empty response for AAAA requests of chinese websites

How-to
----

- Read [CoreDNS Documentation](https://coredns.io/manual/toc/)
- `go get github.com/coredns/coredns`
- `go get github.com/blahgeek/coredns-chinadns`
- Edit `coredns/plugin.cfg`, add `chinadns:github.com/blahgeek/coredns-chinadns` before the hosts middleware (the execution order of plugins are determined by the order of this list)
- `go generate && go build`, test the build using `./coredns -plugins | grep chinadns`

Configuration
----


```

.:5300 {
    log
    # Perform lookups for chinese websites using 114.114.114.114
    # Same parameters as the [proxy plugin](https://coredns.io/plugins/proxy) can be used
    chinadns /etc/coredns/GeoLite2-Country.mmdb . 114.114.114.114:53
    # Use 1.1.1.1 to perform lookups for non-chinese websites
    proxy . 1.1.1.1:53
}

```
