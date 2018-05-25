package nightwatcher

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"time"

	"sync"

	"strings"

	httpstat "github.com/tcnksm/go-httpstat"
)

type WatchItem struct {
	TargetURL string
	Method    string
	HttpBody  string
	Interval  time.Duration
	TimeOut   time.Duration
}

type Watcher struct {
	WatchItems []WatchItem
	wg         sync.WaitGroup
}

func (this *Watcher) Watch() {
	for _, wi := range this.WatchItems {
		this.wg.Add(1)
		go func(wi WatchItem) {
			timeout := time.After(wi.TimeOut)
			heartbeat := time.Tick(wi.Interval)
			for {
				select {
				case <-timeout:
					this.wg.Done()
					return
				case <-heartbeat:
					method := wi.Method
					if wi.HttpBody != "" {
						method = "POST"
					}
					if method == "" {
						method = "GET"
					}
					err := Look(method, wi.TargetURL, wi.HttpBody)
					if err != nil {
						log.Printf("Look %v failed  due to err: %v", wi.TargetURL, err)
					}
				}
			}
		}(wi)
	}
	this.wg.Wait()
}

func Look(method string, url string, body string) (err error) {
	var aim *Aim
	aim, err = glance(method, url, body)
	if err != nil {
		return
	}
	for _, reporter := range Reporters {
		err = reporter.Output(aim)
		if err != nil {
			return
		}
	}
	return
}

func glance(method string, url string, body string) (r *Aim, err error) {
	var result httpstat.Result
	method = strings.ToUpper(method)

	bodyReader := createBody(body)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return
	}

	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	// TODO remove hardcode, make it configurable
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},

		// NOTE do not follow http redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	if _, err = io.Copy(ioutil.Discard, res.Body); err != nil {
		return
	}

	res.Body.Close()
	bodyReadTime := time.Now()
	result.End(bodyReadTime)
	r = NewAim(result, res.StatusCode, bodyReadTime, url)
	return
}
