package nightwatcher

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"time"

	"sync"

	httpstat "github.com/tcnksm/go-httpstat"
)

type WatchItem struct {
	TargetURL string
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
					err := Look(wi.TargetURL)
					if err != nil {
						log.Printf("Look %v failed  due to err: %v", wi.TargetURL, err)
					}
				}
			}
		}(wi)
	}
	this.wg.Wait()
}

func Look(url string) (err error) {
	var aim *Aim
	aim, err = glance(url)
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

func glance(url string) (r *Aim, err error) {
	var result httpstat.Result
	req, err := http.NewRequest("GET", url, nil)
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
	//client := http.DefaultClient
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
