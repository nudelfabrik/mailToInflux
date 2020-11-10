package mtasts

import (
	"encoding/json"
)

type Report struct {
	OrgName string `json:"organization-name"`
}

func ParseMTASTS(raw []byte) (*Report, error) {
	var report Report

	err := json.Unmarshal(raw, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}
