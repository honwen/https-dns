package gdns

import (
	"github.com/golang/glog"
	"github.com/miekg/dns"
)

// HandlerOptions specifies options to be used when instantiating a handler
type HandlerOptions struct{}

// Handler represents a DNS handler
type Handler struct {
	options  *HandlerOptions
	provider Provider
}

// NewHandler creates a new Handler
func NewHandler(provider Provider, options *HandlerOptions) *Handler {
	return &Handler{options, provider}
}

// Handle handles a DNS request
func (h *Handler) Handle(w dns.ResponseWriter, r *dns.Msg) {
	q := DNSQuestion{
		Name:   r.Question[0].Name,
		Type:   r.Question[0].Qtype,
		Subnet: nil,
	}

	edns := r.IsEdns0()
	if edns != nil {
		for _, opt := range edns.Option {
			if opt.Option() == dns.EDNS0SUBNET {
				subnet, ok := opt.(*dns.EDNS0_SUBNET)
				if ok {
					q.Subnet = subnet
				} else {
					glog.V(LERROR).Infoln("parse edns-client-subnet failed", opt.Option())
				}
			}
		}
	}

	switch q.Type {
	case dns.TypeANY:
		glog.V(LINFO).Infoln("request-Blocked", q.Name, dns.TypeToString[q.Type])
		resp := dns.Msg{
			MsgHdr: dns.MsgHdr{
				Id:       r.Id,
				Response: false,
			},
			Compress: r.Compress,
		}
		// Write the response
		if err := w.WriteMsg(&resp); err != nil {
			glog.V(LERROR).Infoln("provider failed", err)
		}
		return
	default:
		glog.V(LINFO).Infoln("requesting", q.Name, dns.TypeToString[q.Type])
	}

	dnsResp, err := h.provider.Query(q)
	if err != nil {
		glog.V(LERROR).Infoln("provider failed", err)
		dns.HandleFailed(w, r)
		return
	}

	questions := []dns.Question{}
	for idx, c := range dnsResp.Question {
		questions = append(questions, dns.Question{
			Name:   c.Name,
			Qtype:  c.Type,
			Qclass: r.Question[idx].Qclass,
		})
	}

	// Parse google RRs to DNS RRs
	answers := transformRR(dnsResp.Answer, "answer")
	authorities := transformRR(dnsResp.Authority, "authority")
	extras := dnsResp.Extra

	resp := dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id:                 r.Id,
			Response:           (dnsResp.ResponseCode == 0),
			Opcode:             dns.OpcodeQuery,
			Authoritative:      false,
			Truncated:          dnsResp.Truncated,
			RecursionDesired:   dnsResp.RecursionDesired,
			RecursionAvailable: dnsResp.RecursionAvailable,
			AuthenticatedData:  dnsResp.AuthenticatedData,
			CheckingDisabled:   dnsResp.CheckingDisabled,
			Rcode:              dnsResp.ResponseCode,
		},
		Compress: r.Compress,
		Question: questions,
		Answer:   answers,
		Ns:       authorities,
		Extra:    extras,
	}

	// glog.V(LINFO).Infof("%+v", resp)
	// Write the response
	if err = w.WriteMsg(&resp); err != nil {
		glog.V(LERROR).Infoln("provider failed", err)
	}
}

// for a given []DNSRR, transform to dns.RR, logging if any errors occur
func transformRR(rrs []DNSRR, logType string) []dns.RR {
	var t []dns.RR

	for _, r := range rrs {
		if rr, err := r.DNSRR(); err != nil {
			glog.V(LERROR).Infoln("unable to translate record rr", logType, r, err)
		} else {
			t = append(t, rr)
		}
	}

	return t
}
