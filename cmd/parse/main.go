package main

import (
	"compress/gzip"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nudelfabrik/mailToInflux/influx"
	"github.com/nudelfabrik/mailToInflux/mtasts"
	"github.com/nudelfabrik/mailToInflux/settings"
)

func main() {
	stdin := flag.Bool("stdin", false, "Use Standard In for file")

	flag.Parse()

	var file *os.File

	settings, err := settings.LoadSettings("")

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
		db, err := influx.NewInfluxDB(settings)
		if err != nil {
			log.Println(err)
		}
		writeAPI := db.Client.WriteAPIBlocking(settings.Org, settings.Bucket)
		p := influxdb2.NewPointWithMeasurement("mtasts").
			AddTag("orgname", report.OrgName).
			AddField("success", report.Success).
			AddField("failure", report.Failure).
			SetTime(report.EndTime)
		err = writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			log.Println(err)
		}
	}
}

func unGz(file *os.File) ([]byte, error) {
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(gzReader)
}
