package mtasts

import (
	"encoding/json"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Report struct {
	OrgName   string
	BeginTime time.Time
	EndTime   time.Time
	Success   int
	Failure   int
}

type rawReport struct {
	OrgName   string   `json:"organization-name"`
	Daterange drange   `json:"date-range"`
	Policies  []policy `json:"policies"`
}

type drange struct {
	BeginTime MtaTime `json:"start-datetime"`
	EndTime   MtaTime `json:"end-datetime"`
}

type policy struct {
	S summary `json:"summary"`
}

type summary struct {
	Success int `json:"total-successful-session-count"`
	Failure int `json:"total-failure-session-count"`
}

func Parse(raw []byte) (*Report, error) {
	var rReport rawReport

	err := json.Unmarshal(raw, &rReport)
	if err != nil {
		return nil, err
	}

	r := Report{}
	r.OrgName = rReport.OrgName
	r.BeginTime = rReport.Daterange.BeginTime.Time
	r.EndTime = rReport.Daterange.EndTime.Time

	for _, p := range rReport.Policies {
		r.Success += p.S.Success
		r.Failure += p.S.Failure
	}

	return &r, nil
}

func (r Report) Measurements() (points []*write.Point) {
	p := influxdb2.NewPointWithMeasurement("mtasts").
		AddTag("orgname", r.OrgName).
		AddField("success", r.Success).
		AddField("failure", r.Failure).
		SetTime(r.EndTime)

	points = append(points, p)

	return points
}
