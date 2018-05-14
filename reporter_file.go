package nightwatcher

import (
	"fmt"
	"os"
)

func NewFileReporter(filename string) (r *FileReporter, err error) {
	var f *os.File
	f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0600)
	if err != nil {
		f.Close()
		return
	}
	r = &FileReporter{f}
	return
}

type FileReporter struct {
	file *os.File
}

func (this *FileReporter) Output(aim *Aim) (err error) {
	var line string
	bodyReadTime := aim.BodyReadTime
	line = fmt.Sprintf(
		"[%v] from %v ping %v, get status: %v, DNSLookup: %v, TCPConnection: %v, TLSHandshake: %v, ServerProcessing: %v ContentTransfer: %v, nameLookup: %v, connect: %v pretransfer: %v, starttransfer: %v, total: %v\n",
		bodyReadTime.Format(nightWatcherLogTimeFormat),
		aim.FromIP(),
		aim.TargetURL,
		aim.StatusCode,
		aim.Result.DNSLookup.String(),
		aim.Result.TCPConnection.String(),
		aim.Result.TLSHandshake.String(),
		aim.Result.ServerProcessing.String(),
		aim.Result.ContentTransfer(bodyReadTime).String(),
		aim.Result.NameLookup.String(),
		aim.Result.Connect.String(),
		aim.Result.Pretransfer.String(),
		aim.Result.StartTransfer.String(),
		aim.Result.Total(bodyReadTime).String(),
	)
	this.file.WriteString(line)
	return
}
