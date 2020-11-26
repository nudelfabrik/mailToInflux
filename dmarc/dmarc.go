package dmarc

import (
	"encoding/xml"
	"log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Report struct {
	OrgName   string   `xml:"report_metadata>org_name"`
	BeginTime UnixTime `xml:"report_metadata>date_range>begin"`
	EndTime   UnixTime `xml:"report_metadata>date_range>end"`
	Records   []Record `xml:"record"`
}

type Record struct {
	SourceIP   string        `xml:"row>source_ip"`
	Count      string        `xml:"row>count"`
	DkimResult PolicySuccess `xml:"row>policy_evaluated>dkim"`
	SpfResult  PolicySuccess `xml:"row>policy_evaluated>spf"`
}

func Parse(raw []byte) (*Report, error) {
	report := Report{}

	err := xml.Unmarshal(raw, &report)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &report, nil
}

func (r Report) Measurement() *write.Point {
	p := influxdb2.NewPointWithMeasurement("dmarc").
		AddTag("orgname", r.OrgName).
		SetTime(r.EndTime.Time)

	return p
}
