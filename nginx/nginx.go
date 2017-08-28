package nginx

import (
	"fmt"
	"io"
	"strings"
	"time"
	"github.com/satyrius/gonx"

	"postdata-replay/helpers"
)

const (
	nginxTimeLayout = "2/Jan/2006:15:04:05 -0700"
)

// NginxReader implements LogReader interface
type NginxReader struct {
	GonxReader *gonx.Reader
}

func parseRequest(requestString string) ([]string, error) {
	parsedRequest := strings.SplitN(requestString, " ", 3)

	if len(parsedRequest) != 3 {
		return parsedRequest, fmt.Errorf("ERROR while parsing string: %s", requestString)
	}

	return parsedRequest, nil
}

func parseNginxTime(timeLocal string) time.Time {
	t, err := time.Parse(nginxTimeLayout, timeLocal)

	helpers.CheckErr(err)

	return t
}

// NewNginxReader creates new reader for a haproxy log format using provided io.Reader
func NewNginxReader(inputReader io.Reader, format string) helpers.LogReader {
	var reader NginxReader
	reader.GonxReader = gonx.NewReader(inputReader, format)

	return &reader
}

func (r *NginxReader) Read() (*helpers.LogEntry, error) {
	var entry helpers.LogEntry

	rec, err := r.GonxReader.Read()

	if err != nil {
		return &entry, err
	}

	timeLocal, err := rec.Field("time_local")

	if err != nil {
		return &entry, err
	}

	requestString, err := rec.Field("request")

	if err != nil {
		return &entry, err
	}

	postData, err := rec.Field("payload")

	if err != nil {
		return &entry, err
	}

	parsedRequest, err := parseRequest(requestString)

	if err != nil {
		return &entry, err
	}

	entry.Method = parsedRequest[0]
	entry.URL = parsedRequest[1]
	entry.Time = parseNginxTime(timeLocal)
	entry.Payload = postData

	return &entry, nil
}