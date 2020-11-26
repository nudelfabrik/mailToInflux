package mailtoinflux

import (
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Report interface {
	Measurements() []*write.Point
}
