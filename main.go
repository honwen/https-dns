package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chenhw2/google-https-dns/gdns"
	"github.com/golang/glog"
	"github.com/miekg/dns"
	"github.com/urfave/cli"
)

var (
	version = "MISSING build version [git hash]"

	gdnsOPT   gdns.GDNSOptions
	gdnsEndPT string

	listenAddress   string
	listenProtocols []string
)

func serve(net, addr string) {
	glog.V(LINFO).Infof("starting %s service on %s", net, addr)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	server := &dns.Server{Addr: addr, Net: net, TsigSecret: nil}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			glog.Errorf("Failed to setup the %s server: %s\n", net, err.Error())
			sig <- syscall.SIGTERM
		}
	}()

	// serve until exit
	<-sig

	glog.V(LINFO).Infof("shutting down %s on interrupt\n", net)
	if err := server.Shutdown(); err != nil {
		glog.V(LERROR).Infof("got unexpected error %s", err.Error())
	}
}

func init() {
	// seed the global random number generator, used in some utilities and the
	// google provider
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	app := cli.NewApp()
	app.Name = "google-https-dns"
	app.Usage = "A DNS-protocol proxy for Google's DNS-over-HTTPS service."
	app.Version = version
	// app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Value: ":5300",
			Usage: "Serve address",
		},
		cli.StringFlag{
			Name:  "endpoint",
			Value: "https://dns.google.com/resolve",
			Usage: "Google DNS-over-HTTPS endpoint url",
		},
		cli.StringFlag{
			Name:  "endpoint-ips",
			Usage: "IPs of the Google DNS-over-HTTPS endpoint; if provided, endpoint lookup skip",
		},
		cli.StringFlag{
			Name:  "dns-servers, d",
			Usage: "DNS Servers used to look up the endpoint; system default is used if absent.",
		},
		cli.StringFlag{
			Name:  "edns, e",
			Usage: "Extension mechanisms for DNS (EDNS) is parameters of the Domain Name System (DNS) protocol.",
		},
		cli.BoolFlag{
			Name:  "no-pad",
			Usage: "Disable padding of Google DNS-over-HTTPS requests to identical length",
		},
		cli.BoolFlag{
			Name:  "udp, U",
			Usage: "Listen on UDP",
		},
		cli.BoolFlag{
			Name:  "tcp, T",
			Usage: "Listen on TCP",
		},
	}
	app.Action = func(c *cli.Context) error {
		glogGangstaShim(c)
		listenAddress = c.String("listen")
		gdnsEndPT = c.String("endpoint")
		if c.Bool("tcp") {
			listenProtocols = append(listenProtocols, "tcp")
		}
		if c.Bool("udp") {
			listenProtocols = append(listenProtocols, "udp")
		}
		if 0 == len(listenProtocols) {
			cli.ShowAppHelp(c)
			os.Exit(0)
		}
		endPtIPs, err := gdns.CSVtoIPs(c.String("endpoint-ips"))
		if err != nil {
			glog.V(LFATAL).Infof("error parsing endpoint-ips: %v", err)
		}
		dnsIPs, err := gdns.CSVtoEndpoints(c.String("dns-servers"))
		if err != nil {
			glog.V(LFATAL).Infof("error parsing dns-servers: %v", err)
		}
		gdnsOPT.EDNS = c.String("edns")
		gdnsOPT.Pad = !c.Bool("no-pad")
		gdnsOPT.DNSServers = dnsIPs
		gdnsOPT.EndpointIPs = endPtIPs
		return nil
	}
	app.Flags = append(app.Flags, glogGangstaFlags...)
	app.Run(os.Args)
	defer glog.Flush()

	provider, err := gdns.NewGDNSProvider(gdnsEndPT, &gdnsOPT)
	if err != nil {
		glog.Exitln(err)
	}
	// options := &gdns.HandlerOptions{}
	handler := gdns.NewHandler(provider, new(gdns.HandlerOptions))
	dns.HandleFunc(".", handler.Handle)

	// start the servers
	servers := make(chan bool)
	for _, protocol := range listenProtocols {
		go func(protocol string) {
			serve(protocol, listenAddress)
			servers <- true
		}(protocol)
	}

	// wait for servers to exit
	for range listenProtocols {
		<-servers
	}

	glog.V(LINFO).Infoln("servers exited, stopping")
}
