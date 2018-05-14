package nightwatcher

import (
	"fmt"

	"github.com/fatih/color"
)

var DefaultReporter = StdoutReporter{}

type StdoutReporter struct {
}

func (this StdoutReporter) Output(aim *Aim) (err error) {
	fmt.Printf("[%v] watch %v status: %v, cost: %v\n",
		aim.BodyReadTime.Format(nightWatcherLogTimeFormat),
		color.BlueString(aim.TargetURL),
		colorizeStatusCode(aim.StatusCode),
		fmtb(aim.Result.Total(aim.BodyReadTime)),
	)

	result := aim.Result
	bodyReadTime := aim.BodyReadTime
	printf(colorize(HTTPSTemplate),
		fmta(result.DNSLookup),                     // dns lookup
		fmta(result.Connect),                       // tcp connection
		fmta(result.TLSHandshake),                  // tls handshake
		fmta(result.ServerProcessing),              // server processing
		fmta(result.ContentTransfer(bodyReadTime)), // content transfer
		fmtb(result.NameLookup),                    // namelookup
		fmtb(result.Connect),                       // connect
		fmtb(result.Pretransfer),                   // pretransfer
		fmtb(result.StartTransfer),                 // starttransfer
		fmtb(result.Total(bodyReadTime)),           // total
	)
	return
}
