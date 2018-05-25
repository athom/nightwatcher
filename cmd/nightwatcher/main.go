package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/athom/nightwatcher"
	"github.com/jinzhu/configor"
)

var (
	showVersion        bool
	silence            bool
	outputFile         string
	configFile         string
	httpMethod         string
	httpPostBody       string
	elasticsearchURL   string
	elasticsearchIndex string

	interval int
	duration int

	version = "devel"
)

func init() {
	flag.IntVar(&interval, "i", 3, "interval seconds for watch period")
	flag.IntVar(&duration, "d", 4, "duration of watch time")
	flag.BoolVar(&silence, "s", false, "make stdout silence")
	flag.StringVar(&httpMethod, "X", "GET", "use config file")
	flag.StringVar(&httpPostBody, "D", "", "the body of a POST or PUT request; from file use @filename")
	flag.StringVar(&configFile, "f", "", "use config file")
	flag.StringVar(&outputFile, "o", "", "use output file to report watch result")
	flag.StringVar(&elasticsearchURL, "e", "", "use elasticsearch to collect watch result")
	flag.StringVar(&elasticsearchIndex, "x", "", "index for the collected result in elasticsearch")
	flag.BoolVar(&showVersion, "v", false, "print version number")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] URL\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
}

func runWithConfig(configFile string) (err error) {
	type WatchItem struct {
		TargetURL string        `yaml:"url"`
		Method    string        `yaml:"method"`
		HttpBody  string        `yaml:"http_body"`
		Interval  time.Duration `yaml:"interval"`
		Duration  time.Duration `yaml:"duration"`
	}

	type StdoutReporter struct {
		Enabled bool `yaml:"enabled"`
	}

	type FileReporter struct {
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"path"`
	}

	type ElasticsearchReporter struct {
		Enabled bool   `yaml:"enabled"`
		URL     string `yaml:"url"`
		Index   string `yaml:"index"`
	}

	type Reporters struct {
		StdoutReporter        StdoutReporter        `yaml:"stdout"`
		FileReporter          FileReporter          `yaml:"file"`
		ElasticsearchReporter ElasticsearchReporter `yaml:"elasticsearch"`
	}

	type NightWatcherConfig struct {
		WatchItems []WatchItem `yaml:"watch_items"`
		Reporters  Reporters   `yaml:"reporters"`
	}

	var conf NightWatcherConfig
	err = configor.Load(&conf, configFile)
	if err != nil {
		return err
	}

	var watcherItems []nightwatcher.WatchItem
	for _, item := range conf.WatchItems {
		watcherItems = append(watcherItems, nightwatcher.WatchItem{
			TargetURL: item.TargetURL,
			Method:    item.Method,
			HttpBody:  item.HttpBody,
			Interval:  item.Interval,
			TimeOut:   item.Duration,
		})
	}
	watcher := nightwatcher.Watcher{WatchItems: watcherItems}

	if !conf.Reporters.StdoutReporter.Enabled {
		nightwatcher.ClearReporters()
	}
	if conf.Reporters.FileReporter.Enabled {
		path := conf.Reporters.FileReporter.Path
		reporter, err := nightwatcher.NewFileReporter(path)
		if err != nil {
			panic(err)
		}
		nightwatcher.AddReporter(reporter)
	}
	if conf.Reporters.ElasticsearchReporter.Enabled {
		elasticsearchIndex := conf.Reporters.ElasticsearchReporter.Index
		elasticsearchURL := conf.Reporters.ElasticsearchReporter.URL
		if elasticsearchIndex == "" {
			elasticsearchIndex = "nightwatcher"
		}
		esReporter, err := nightwatcher.NewElasticsearchReporter(
			[]string{elasticsearchURL},
			elasticsearchIndex,
		)
		if err != nil {
			panic(err)
		}
		nightwatcher.AddReporter(esReporter)
	}

	if len(nightwatcher.Reporters) == 0 {
		fmt.Printf("%s must spesify at least one reporter to handle the result, such as stdout, output file, or elasticsearch\n", os.Args[0])
		os.Exit(1)
	}

	watcher.Watch()
	return
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Printf("%s %s (runtime: %s)\n", os.Args[0], version, runtime.Version())
		os.Exit(0)
	}

	if configFile != "" {
		runWithConfig(configFile)
		return
	}

	n := len(os.Args)
	if n <= 1 {
		usage()
		return
	}

	url := os.Args[n-1]

	if (httpMethod == "POST" || httpMethod == "PUT") && httpPostBody == "" {
		log.Fatal("must supply post body using -d when POST or PUT is used")
	}

	watcher := nightwatcher.Watcher{
		WatchItems: []nightwatcher.WatchItem{
			{
				Method:    httpMethod,
				HttpBody:  httpPostBody,
				TargetURL: url,
				TimeOut:   time.Duration(duration) * time.Second,
				Interval:  time.Duration(interval) * time.Second,
			},
		},
	}

	if silence {
		nightwatcher.ClearReporters()
	}

	if outputFile != "" {
		reporter, err := nightwatcher.NewFileReporter(outputFile)
		if err != nil {
			panic(err)
		}
		nightwatcher.AddReporter(reporter)
	}

	if elasticsearchURL != "" {
		if elasticsearchIndex == "" {
			elasticsearchIndex = "nightwatcher"
		}
		esReporter, err := nightwatcher.NewElasticsearchReporter(
			[]string{elasticsearchURL},
			elasticsearchIndex,
		)
		if err != nil {
			panic(err)
		}
		nightwatcher.AddReporter(esReporter)
	}

	if len(nightwatcher.Reporters) == 0 {
		fmt.Printf("%s must spesify at least one reporter to handle the result, such as stdout, output file, or elasticsearch\n", os.Args[0])
		os.Exit(1)
	}

	watcher.Watch()
	return
}
