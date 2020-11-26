package influx

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	mti "github.com/nudelfabrik/mailToInflux"
	"github.com/nudelfabrik/mailToInflux/settings"
)

type DB struct {
	Client influxdb2.Client
	api    api.WriteAPIBlocking
}

func NewInfluxDB(settings *settings.Settings) (*DB, error) {
	db := DB{}

	db.Client = influxdb2.NewClient(settings.URL, settings.Token)

	db.api = db.Client.WriteAPIBlocking(settings.Org, settings.Bucket)

	return &db, nil
}

func (db DB) Write(ctx context.Context, r mti.Report) error {
	for _, m := range r.Measurements() {
		err := db.api.WritePoint(ctx, m)
		if err != nil {
			return err
		}
	}
	return nil
}
