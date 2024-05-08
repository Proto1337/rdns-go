package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "net"
    "strings"
    "os"
)

var (
    DNS_SERVER      string
    DNS_PORT        uint
    DNS_PROTOCOL    string
    IP              string
)

func setFlags() {
    /*
    * Set the values of the variables based on commandline flags
    */
    flag.StringVar(&DNS_SERVER, "r", "", "Which DNS server to use as resolver. Specify IP:PORT.")
    flag.StringVar(&DNS_PROTOCOL, "p", "udp", "Which protocol to use when contacting the DNS server.")
    flag.StringVar(&IP, "i", "127.0.0.1", "Which IP to resolve.")
    flag.Parse()
}

func resolveIP(ctx context.Context, ip string, dns_resolver string, dns_resolver_protocol string) string {
    var resolver *net.Resolver
    if dns_resolver != "" {
        resolver = &net.Resolver{
            PreferGo: true,
            Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
                d := net.Dialer{}
                return d.DialContext(ctx, dns_resolver_protocol, dns_resolver)
            },
        }
    }
    res_addr, err := resolver.LookupAddr(ctx, ip)
    if err != nil {
        log.Fatalf("Error resolving given IP address.")
        os.Exit(1)
    }
    addr := strings.TrimRight(res_addr[0], ".")
    return addr
}

func main() {
    ctx := context.Background()
    setFlags()
    addr := resolveIP(ctx, IP, DNS_SERVER, DNS_PROTOCOL)
    fmt.Println(addr)
}
