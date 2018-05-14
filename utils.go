package nightwatcher

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// shame! these handy methods are private, have to stole from https://github.com/davecheney/httpstat
const (
	HTTPSTemplate = `` +
		`  DNS Lookup   TCP Connection   TLS Handshake   Server Processing   Content Transfer` + "\n" +
		`[%s  |     %s  |    %s  |        %s  |       %s  ]` + "\n" +
		`            |                |               |                   |                  |` + "\n" +
		`   namelookup:%s      |               |                   |                  |` + "\n" +
		`                       connect:%s     |                   |                  |` + "\n" +
		`                                   pretransfer:%s         |                  |` + "\n" +
		`                                                     starttransfer:%s        |` + "\n" +
		`                                                                                total:%s` + "\n"

	HTTPTemplate = `` +
		`   DNS Lookup   TCP Connection   Server Processing   Content Transfer` + "\n" +
		`[ %s  |     %s  |        %s  |       %s  ]` + "\n" +
		`             |                |                   |                  |` + "\n" +
		`    namelookup:%s      |                   |                  |` + "\n" +
		`                        connect:%s         |                  |` + "\n" +
		`                                      starttransfer:%s        |` + "\n" +
		`                                                                 total:%s` + "\n"
)

func grayscale(code color.Attribute) func(string, ...interface{}) string {
	return color.New(code + 232).SprintfFunc()
}

func printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(color.Output, format, a...)
}

func fmta(d time.Duration) string {
	return color.CyanString("%7dms", int(d/time.Millisecond))
}

func fmtb(d time.Duration) string {
	return color.CyanString("%-9s", strconv.Itoa(int(d/time.Millisecond))+"ms")
}

func colorize(s string) string {
	v := strings.Split(s, "\n")
	v[0] = grayscale(16)(v[0])
	return strings.Join(v, "\n")
}

func colorizeStatusCode(code int) (r string) {
	if code >= 400 {
		return color.RedString("%v", code)
	}
	if code >= 300 {
		return color.YellowString("%v", code)
	}
	return color.GreenString("%v", code)
}

func getLocalAddress() (r *net.UDPAddr) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	r = conn.LocalAddr().(*net.UDPAddr)
	return
}
