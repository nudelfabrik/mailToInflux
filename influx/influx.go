package influx

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nudelfabrik/mailToInflux/settings"
)

type DB struct {
	Client influxdb2.Client
}

func NewInfluxDB(settings *settings.Settings) (*DB, error) {
	db := DB{}

	client := influxdb2.NewClient(settings.URL, settings.Token)

	db.Client = client

	return &db, nil
}
