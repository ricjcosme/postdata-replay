package main

import (
	// "fmt"
	"io"
	"os"
	"flag"
	"log"
	"strings"
	"compress/gzip"

	"github.com/ricjcosme/postdata-replay/nginx"
	"github.com/ricjcosme/postdata-replay/helpers"
)

func main() {

	log.Printf("Started\n")

	flag.Parse()

	if helpers.Debug {
		log.Printf("Parsing %s log file\n", helpers.InputLogFile)
	}

	var inputReader io.Reader
	file, err := os.Open(helpers.InputLogFile)

	defer file.Close()

	if strings.HasSuffix(helpers.InputLogFile, "gz") {
		inputReader, err = gzip.NewReader(file)
		helpers.CheckErr(err)
	} else {
		inputReader = file
	}

	var reader helpers.LogReader

	reader = nginx.NewNginxReader(inputReader, helpers.Format)

	helpers.ReadLog(reader)

	log.Printf("Finished\n")
}
