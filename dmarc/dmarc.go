package dmarc

import (
	"encoding/xml"
	"log"
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

func ParseDmarc(raw []byte) (*Report, error) {
	var report *Report

	err := xml.Unmarshal(raw, report)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return report, nil
}
