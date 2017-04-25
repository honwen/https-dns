package gdns

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// CSVtoEndpoints takes a comma-separated string of endpoints, and parses to a
// []gdns.Endpoint
func CSVtoEndpoints(csv string) (eps []Endpoint, err error) {
	reps := strings.Split(csv, ",")
	for _, r := range reps {
		if r == "" {
			continue
		}

		ep, err := ParseEndpoint(r, 53)
		if err != nil {
			return eps, err
		}

		eps = append(eps, ep)
	}

	return eps, err
}

// CSVtoIPs takes a comma-separated string of IPs, and parses to a []net.IP
func CSVtoIPs(csv string) (ips []net.IP, err error) {
	rs := strings.Split(csv, ",")

	for _, r := range rs {
		if r == "" {
			continue
		}

		ip := net.ParseIP(r)
		if ip == nil {
			return ips, fmt.Errorf("unable to parse IP from string %s", r)
		}
		ips = append(ips, ip)
	}

	return
}
