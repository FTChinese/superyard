package subs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FTChinese/go-rest/enum"
)

type Kind int

const (
	KindDeny Kind = iota // If user remaining subs period exceed max allowable one, or any other error.
	KindCreate
	KindRenew
	KindUpgrade
)

var subsKindNames = [...]string{
	"",
	"create",
	"renew",
	"upgrade",
}

var subsKindMap = map[Kind]string{
	1: subsKindNames[1],
	2: subsKindNames[2],
	3: subsKindNames[3],
}

var subsKindValue = map[string]Kind{
	subsKindNames[1]: 1,
	subsKindNames[2]: 2,
	subsKindNames[3]: 3,
}

func ParseSubsKind(name string) (Kind, error) {
	if x, ok := subsKindValue[name]; ok {
		return x, nil
	}

	return KindDeny, fmt.Errorf("%s is not valid Kind", name)
}

func (x Kind) String() string {
	if s, ok := subsKindMap[x]; ok {
		return s
	}
	return ""
}

func (x Kind) SnapshotReason() enum.SnapshotReason {
	switch x {
	case KindRenew:
		return enum.SnapshotReasonRenew
	case KindUpgrade:
		return enum.SnapshotReasonUpgrade
	default:
		return enum.SnapshotReasonNull
	}
}

func (x *Kind) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tmp, _ := ParseSubsKind(s)

	*x = tmp

	return nil
}

func (x Kind) MarshalJSON() ([]byte, error) {
	s := x.String()

	if s == "" {
		return []byte("null"), nil
	}

	return []byte(`"` + s + `"`), nil
}

func (x *Kind) Scan(src interface{}) error {
	if src == nil {
		*x = KindDeny
		return nil
	}

	switch s := src.(type) {
	case []byte:
		tmp, _ := ParseSubsKind(string(s))
		*x = tmp
		return nil

	default:
		return errors.New("incompatible type to scan")
	}
}

func (x Kind) Value() (driver.Value, error) {
	s := x.String()
	if s == "" {
		return nil, nil
	}

	return s, nil
}
