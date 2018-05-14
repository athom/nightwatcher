# Night Watcher


Night Watcher is a small go based application that make httpstat to websites easier. 

[![Build Status](https://api.travis-ci.org/athom/nightwatcher.png?branch=master)](https://travis-ci.org/athom/nightwatcher)


## Installation

```bash
	go get "github.com/athom/nightwatcher/cmd/nightwatcher"
```

or 

```bash
git clone git@github.com:athom/nightwatcher.git
cd $GOPATH/github.com/athom/nightwatcher
govender sync
cd $GOPATH/github.com/athom/nightwatcher/cmd/nightwatcher
go install
```

## Usage

- **Ping URL with status code and time consume statistics**

  The simplest usage is designed to replace `httpstat`. Just 

```sh
nightwatcher https://about.gitlab.com
```

- **Watch periodly**

  By default it only ping the URL for once, but you can provide a interval and duration to make it run periodly. 

This example will ping github.com every 2 seconds for 1 minute.  

```sh
nightwatcher -i 2 -d 60 https://about.gitlab.com
```

- **Record the result in local file**

  By providing the output filename, the watch result can be saved for persistence with unified format for further analysis.  
  
```sh
nightwatcher -i 2 -d 60 -o watch_result.txt https://about.gitlab.com
```

- **Track the records to Elasticsearch**

  Human friendly is one important factor when designing nightwatcher. So the result is well structed and can be easily store to storages. By default we have local file storage and remote Elasticsearch. 
  Here is the usage of how to spy/monitor website and them store the result to Elasticsearch.  
  
```sh
nightwatcher -i 2 -d 60 -e your_fancy_es_url https://about.gitlab.com
```

- **Support Multiple URLs & configure file**
  The demand of watch multiple websites is very common, you can still use nightwatcher to do that for you by specified a yaml config.  

```sh
nightwatcher -f my_fancy_watcher_config.yml
```

Example config file could be look like:

```yaml
watch_items:
- url: https://github.com
  interval: 3s
  duration: 5m

- url: http://github.com
  interval: 3s
  duration: 5m

- url: https://cloud.tencent.com
  interval: 3s
  duration: 5m

- url: https://about.gitlab.com
  interval: 3s
  duration: 5m

reporters:
  stdout:
    enabled: yes
  file:
    enabled: yes
    path: "./nightwatcher_result.txt"
  elasticsearch:
    enabled: yes
    url: "https://f6swlx6x1k:v8b332xz70@nightwatcher-storage-5963183380.us-east-1.bonsaisearch.net"
    index: "nightwatcher"
```



## Plan & Todo

- Support more http methods.
- Support more storage plugins.

## License

Night Watcher is released under the [WTFPL License](http://www.wtfpl.net/txt/copying).

