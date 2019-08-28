package reader

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type SnapshotReason int

const (
	SnapshotReasonNull SnapshotReason = iota
	SnapshotReasonRenew
	SnapshotReasonUpgrade
	SnapshotReasonDelete
)

var snapshotNames = [...]string{
	"",
	"renew",
	"upgrade",
	"delete",
}

// Maps SnapshotReason to string representation
var snapshotToStrings = map[SnapshotReason]string{
	1: snapshotNames[1],
	2: snapshotNames[2],
	3: snapshotNames[3],
}

// Maps a string value to a Snapshot instance.
var snapshotValues = map[string]SnapshotReason{
	snapshotNames[1]: 1,
	snapshotNames[2]: 2,
	snapshotNames[3]: 3,
}

func ParseSnapshotReason(name string) (SnapshotReason, error) {
	if x, ok := snapshotValues[name]; ok {
		return x, nil
	}

	return SnapshotReasonNull, fmt.Errorf("%s is not valid SnapshotReason", name)
}

func (x SnapshotReason) String() string {
	if s, ok := snapshotToStrings[x]; ok {
		return s
	}

	return ""
}

func (x *SnapshotReason) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tmp, _ := ParseSnapshotReason(s)

	*x = tmp

	return nil
}

func (x SnapshotReason) MarshalJSON() ([]byte, error) {
	s := x.String()

	if s == "" {
		return []byte("null"), nil
	}

	return []byte(`"` + s + `"`), nil
}

func (x *SnapshotReason) Scan(src interface{}) error {
	if src == nil {
		*x = SnapshotReasonNull
		return nil
	}

	switch s := src.(type) {
	case []byte:
		tmp, _ := ParseSnapshotReason(string(s))
		*x = tmp
		return nil

	default:
		return errors.New("imcompatible type to scan")
	}
}

func (x SnapshotReason) Value() (driver.Value, error) {
	s := x.String()
	if s == "" {
		return nil, nil
	}

	return s, nil
}
