# Replay Nginx logs (GETs and POSTs)
Based on https://github.com/Gonzih/log-replay

## Installation

```
go get -u github.com/ricjcosme/postdata-replay
```

## Usage

```
Usage of postdata-replay:
  -debug
    	Print extra debugging information
  -file string
    	Log file name to read. Read from STDIN if file name is '-' (default "-")
  -format string
    	Nginx log format (extended "$remote_addr - - [$time_local] \"$request\" $status $request_length $body_bytes_sent $request_time $payload\"$t_size\" $read_time $gen_time")
  -prefix string
    	URL prefix to query (default "http://localhost")
  -ratio int
    	Replay speed ratio, higher means faster replay speed (default 1)
  -skip-sleep
    	Skip sleep between http calls based on log timestapms
  -timeout int
    	Request timeout in milliseconds, 0 means no timeout (default 60000)
```

```bash
# Replay access log
postdata-replay --file my-acces.log --debug

# Duplicate traffic on the staging host
tail -f /var/log/acces.log | postdata-replay --prefix http://staging-host --skip-sleep
```

## GET / POST

Methods GET and POST are currently available for Nginx logs.

## License

[MIT](LICENSE)

[license-url]: LICENSE

[license-image]: https://img.shields.io/github/license/mashape/apistatus.svg

[capture]: capture.png
