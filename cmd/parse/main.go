package main

import (
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/nudelfabrik/mailToInflux/dmarc"
	"github.com/nudelfabrik/mailToInflux/influx"
	"github.com/nudelfabrik/mailToInflux/mtasts"
	"github.com/nudelfabrik/mailToInflux/settings"
)

func main() {
	// Parse Flags
	stdin := flag.Bool("stdin", false, "Use Standard In for file")

	flag.Parse()

	// load influx connection
	settings, err := settings.LoadSettings("")
	if err != nil {
		log.Fatal(err)
	}

	db, err := influx.NewInfluxDB(settings)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Client.Close()

	// Open Reader
	var file *os.File

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

	// try Gzip fist
	report, err := unGz(file)
	if err == nil {
		// When error is nil, ungzip successful, write datapoint and exit
		err = db.Write(context.Background(), report)
		if err != nil {
			log.Println(err)
		}
		return

	} else if !errors.Is(err, gzip.ErrHeader) {
		// If error something else than unsuccessful ungzip, something else is wrong
		// e.g. File error, unmarshaling error. Log error and exit
		log.Println(err)
		return
	}

	dreport, err := unZip(file)

	if err == nil {
		// When error is nil, unzip successful, write datapoint and exit
		err = db.Write(context.Background(), dreport)
		if err != nil {
			log.Println(err)
		}
		return
	}

}

func unZip(zfile *os.File) (*dmarc.Report, error) {
	stats, err := zfile.Stat()
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(zfile, stats.Size())
	if err != nil {
		return nil, err
	}

	for _, f := range zipReader.File {
		file, err := f.Open()
		if err != nil {
			return nil, err
		}

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		report, err := dmarc.Parse(bytes)
		if err != nil {
			log.Println(err)
		} else {
			return report, nil
		}
	}

	return nil, errors.New("No xml file found")

}

func unGz(file io.Reader) (*mtasts.Report, error) {
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	defer gzReader.Close()

	bytes, err := ioutil.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}

	return mtasts.Parse(bytes)
}
