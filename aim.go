package nightwatcher

import (
	"os"
	"time"

	httpstat "github.com/tcnksm/go-httpstat"
)

func NewAim(result httpstat.Result, statusCode int, bodyReadTime time.Time, targetURL string) (r *Aim) {
	var err error
	r = &Aim{Result: result}
	r.StatusCode = statusCode
	r.TargetURL = targetURL
	r.HostName, err = os.Hostname()
	r.BodyReadTime = bodyReadTime
	if err != nil {
		panic(err)
	}
	return r
}

type Aim struct {
	httpstat.Result
	TargetURL    string
	HostName     string
	StatusCode   int
	localIP      string
	BodyReadTime time.Time
}

func (this *Aim) FromIP() (r string) {
	if this.localIP != "" {
		return this.localIP
	}

	localAddr := getLocalAddress()
	this.localIP = localAddr.IP.String()
	return this.localIP
}

func (this *Aim) Durations() map[string]time.Duration {
	r := this.Result
	return map[string]time.Duration{
		"DNSLookup":        r.DNSLookup,
		"TCPConnection":    r.TCPConnection,
		"TLSHandshake":     r.TLSHandshake,
		"ServerProcessing": r.ServerProcessing,
		"ContentTransfer":  r.ContentTransfer(this.BodyReadTime),
		"NameLookup":       r.NameLookup,
		"Connect":          r.Connect,
		"Pretransfer":      r.Connect,
		"StartTransfer":    r.StartTransfer,
		"Total":            r.Total(this.BodyReadTime),
	}
}

func (this *Aim) ToJson() map[string]interface{} {
	r := this.Result
	return map[string]interface{}{
		"HostName":   this.HostName,
		"FromIp":     this.FromIP(),
		"TargetURL":  this.TargetURL,
		"StatusCode": this.StatusCode,

		"BodyReadTime":     this.BodyReadTime,
		"DNSLookup":        r.DNSLookup,
		"TCPConnection":    r.TCPConnection,
		"TLSHandshake":     r.TLSHandshake,
		"ServerProcessing": r.ServerProcessing,
		"ContentTransfer":  r.ContentTransfer(this.BodyReadTime),
		"NameLookup":       r.NameLookup,
		"Connect":          r.Connect,
		"Pretransfer":      r.Connect,
		"StartTransfer":    r.StartTransfer,
		"Total":            r.Total(this.BodyReadTime),
	}
}
