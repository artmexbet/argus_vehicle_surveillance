package models

import (
	"fmt"
	"github.com/google/uuid"
)

type SecurityCarIDType uuid.UUID

func (scid *SecurityCarIDType) UnmarshalJSON(b []byte) error {
	id, err := uuid.Parse(string(b[:]))
	if err != nil {
		return err
	}
	*scid = SecurityCarIDType(id)
	return nil
}

func (scid *SecurityCarIDType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", uuid.UUID(*scid).String())), nil
}
