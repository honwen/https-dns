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
		Name: r.Question[0].Name,
		Type: r.Question[0].Qtype,
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
	}

	// Parse google RRs to DNS RRs
	for _, a := range dnsResp.Answer {
		rr, err := a.RR()
		if err != nil {
			glog.V(LERROR).Infof("%+v", err)
		} else {
			resp.Answer = append(resp.Answer, rr)
		}
	}

	// Parse google RRs to DNS RRs
	for _, ns := range dnsResp.Authority {
		rr, err := ns.RR()
		if err != nil {
			glog.V(LERROR).Infof("%+v", err)
		} else {
			resp.Ns = append(resp.Answer, rr)
		}
	}

	// Parse google RRs to DNS RRs
	for _, extra := range dnsResp.Extra {
		rr, err := extra.RR()
		if err != nil {
			glog.V(LERROR).Infof("%+v", err)
		} else {
			resp.Extra = append(resp.Answer, rr)
		}
	}

	// glog.V(LINFO).Infof("%+v", resp)
	// Write the response
	if err = w.WriteMsg(&resp); err != nil {
		glog.V(LERROR).Infoln("provider failed", err)
	}
}
