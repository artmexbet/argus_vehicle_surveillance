package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"time"
)

type AccountIDType uuid.UUID

func (ait *AccountIDType) UnmarshalJSON(b []byte) error {
	id, err := uuid.Parse(string(b[:]))
	if err != nil {
		return err
	}
	*ait = AccountIDType(id)
	return nil
}

func (ait *AccountIDType) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(*ait).String())
}

type CarIDType int
type TimestampType time.Time

func (t *TimestampType) UnmarshalJSON(b []byte) error {
	tt, err := time.Parse(time.RFC3339, strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}
	*t = TimestampType(tt)
	return nil
}

func (t *TimestampType) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*t).Format(time.RFC3339))
}
