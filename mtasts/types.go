package mtasts

import (
	"encoding/json"
	"time"
)

type MtaTime struct {
	time.Time
}

func (t *MtaTime) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	rawTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	*t = MtaTime{Time: rawTime}

	return nil
}
