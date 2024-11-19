package models

import (
	"fmt"
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
	return []byte(fmt.Sprintf("\"%s\"", uuid.UUID(*ait).String())), nil
}

type CarIDType int
type TimestampType time.Time

func (t *TimestampType) UnmarshalJSON(b []byte) error {
	fmt.Println(strings.Trim(string(b[:]), "\""))
	tt, err := time.Parse(time.RFC3339, strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}
	*t = TimestampType(tt)
	return nil
}

func (t *TimestampType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(*t).Format(time.RFC3339))), nil
}
