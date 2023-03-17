package conv

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
)

type HexStr string

// Scan implements sql.Scanner interface to retrieve binary value from SQL to a hex string.
func (h *HexStr) Scan(src interface{}) error {
	if src == nil {
		*h = ""
		return nil
	}

	switch s := src.(type) {
	case []byte:
		tmp := hex.EncodeToString(s)
		*h = HexStr(tmp)
		return nil

	default:
		return errors.New("incompatible data type to scan")
	}
}

func (h HexStr) Value() (driver.Value, error) {
	b, err := hex.DecodeString(string(h))
	if err != nil {
		return nil, err
	}

	return b, nil
}
