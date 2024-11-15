package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type BUID uuid.UUID

func (u BUID) Value() (driver.Value, error) {
	return uuid.UUID(u).MarshalBinary() // nolint:wrapcheck
}

func (u *BUID) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan UUID: %v", value)
	}

	parsedUUID, err := uuid.FromBytes(bytes)
	if err != nil {
		return err // nolint:wrapcheck
	}

	*u = BUID(parsedUUID)
	return nil
}
