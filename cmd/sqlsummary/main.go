package main

import (
	"flag"
	"log"
	"os"

	"github.com/progfay/sqlsummary"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func main() {
	var filePath = flag.String("f", "", "sql file path (required)")
	var bufferSize = flag.Int("b", 0, "buffer size (default: 1024 * 1024 * 1024 = 1GB)")
	flag.Parse()

	if *bufferSize == 0 {
		*bufferSize = 1024 * 1024 * 1024
	}

	f, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sqlsummary.Run(os.Stdout, f, *bufferSize)
}
