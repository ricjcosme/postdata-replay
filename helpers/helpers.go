package helpers

import (
	"log"
	"time"
	"io"
	"net/http"
	"bytes"
	"flag"
	"sync"
)

// LogEntry is single parsed entry from the log file
type LogEntry struct {
	Time    time.Time
	Method  string
	URL     string
	Payload string
}

// LogReader provides generic log parser interface
type LogReader interface {
	Read() (*LogEntry, error)
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var Format string
var InputLogFile string
var prefix string
var ratio int64
var Debug bool
var clientTimeout int64
var skipSleep bool

var messages chan int
var wg sync.WaitGroup

func init() {
	flag.StringVar(&Format, "format", `$remote_addr - - [$time_local] "$request" $status $request_length $body_bytes_sent $request_time $payload"$t_size" $read_time $gen_time`, "Nginx log format")
	flag.StringVar(&InputLogFile, "file", "-", "Log file name to read. Read from STDIN if file name is '-'")
	flag.StringVar(&prefix, "prefix", "http://localhost:8080", "URL prefix to query")
	flag.Int64Var(&ratio, "ratio", 1, "Replay speed ratio, higher means faster replay speed")
	flag.BoolVar(&Debug, "debug", false, "Print extra debugging information")
	flag.Int64Var(&clientTimeout, "timeout", 60000, "Request timeout in milliseconds, 0 means no timeout")
	flag.BoolVar(&skipSleep, "skip-sleep", false, "Skip sleep between http calls based on log timestapms")

	messages = make(chan int)
}

func ReadLog(reader LogReader) {
	var nilTime time.Time
	var lastTime time.Time

	for {
		rec, err := reader.Read()

		if err == io.EOF {
			log.Println("Reached EOF")
			break
		} else {
			CheckErr(err)
		}

		if !skipSleep {
			if lastTime != nilTime {

				differenceUnix := rec.Time.Sub(lastTime).Nanoseconds()

				if differenceUnix > 0 {
					durationWithRation := time.Duration(differenceUnix / ratio)

					if Debug {
						log.Printf("Sleeping for: %.2f seconds", durationWithRation.Seconds())
					}
					time.Sleep(durationWithRation)
				} else {
					if Debug {
						log.Println("No need for sleep!")
					}
				}

			}

			lastTime = rec.Time
		}

		wg.Add(1)

		go replayHTTPRequest(rec.Method, rec.URL, rec.Payload)

		wg.Wait()

	}
}

func replayHTTPRequest(method string, url string, payload string) {

	defer wg.Done()

	path := prefix + url

	if Debug {
		log.Printf("Querying %s %s %s\n", method, path, payload)
	}

	client := &http.Client{
		Timeout: time.Duration(clientTimeout) * time.Millisecond,
	}

	req, err := http.NewRequest(method, path, bytes.NewBufferString(payload))

	if method == "POST" {
		req.Header.Add("Content-Type", "application/application/octet-stream")
	}

	if err != nil {
		if Debug {
			log.Printf("ERROR %s while creating new request to %s", err, path)
		}
		return
	}

	req.Header.Set("User-Agent", "nginx-postdata-log-replay")

	_, err = client.Do(req)

	if err != nil {
		if Debug {
			log.Printf("ERROR %s in server response", err)
		}
		return
	}

	messages <- 1
}