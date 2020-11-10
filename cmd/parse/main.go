package main

import (
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/nudelfabrik/mailToInflux/mtasts"
)

func main() {
	stdin := flag.Bool("stdin", false, "Use Standard In for file")

	flag.Parse()

	var file *os.File

	var err error

	if *stdin {
		file = os.Stdin
	} else {
		if flag.NArg() == 0 {
			log.Fatal("No filename found.")
		}

		file, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
	}
	defer file.Close()

	if err != nil {
		log.Println(err)
	}

	bytes, err := unGz(file)

	if errors.Is(err, gzip.ErrHeader) {
		log.Println("no gzip file")
	} else {
		report, err := mtasts.ParseMTASTS(bytes)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(report)
	}
}

func unGz(file *os.File) ([]byte, error) {
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(gzReader)
}
