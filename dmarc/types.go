package dmarc

import (
	"encoding/xml"
	"time"
)

type PolicySuccess int

const (
	Pass PolicySuccess = iota
	Fail
	Unknown
)

func (p PolicySuccess) String() string {
	switch p {
	case Pass:
		return "pass"
	case Fail:
		return "fail"
	case Unknown:
		return "Unknown"
	default:
		return ""
	}
}

func (p *PolicySuccess) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string

	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	switch s {
	case "pass":
		*p = Pass
	case "fail":
		*p = Fail
	default:
		*p = Unknown
	}

	return nil
}

type UnixTime struct {
	time.Time
}

func (t *UnixTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v int64

	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}

	timeValue := time.Unix(v, 0)

	*t = UnixTime{Time: timeValue}

	return nil
}
