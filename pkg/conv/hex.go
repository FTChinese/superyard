package conv

import (
	"bytes"
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"errors"
)

// HexBin holds a slice of bytes
// that could be converted to hexdeciaml string,
// used for SQL VARBINARY type.
type HexBin []byte

func RandomHexBin(len int) (HexBin, error) {
	b := make([]byte, len)

	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

// DecodeHexString decodes a hex-encoded string.
func DecodeHexString(s string) (HexBin, error) {
	return hex.DecodeString(s)
}

// DecodeHexBytes decodes a hex-encoded bytes.
func DecodeHexBytes(src []byte) (HexBin, error) {
	n, err := hex.Decode(src, src)
	return src[:n], err
}

func (x HexBin) ToBytes() []byte {
	return []byte(x)
}

func (x HexBin) String() string {
	return hex.EncodeToString(x[:])
}

func (x *HexBin) UnmarshalJSON(b []byte) error {

	if b == nil {
		*x = nil
		return nil
	}

	b = bytes.Trim(b, `"`)

	hb, err := DecodeHexBytes(b)

	if err != nil {
		return err
	}

	*x = hb

	return nil
}

func (x HexBin) MarshalJSON() ([]byte, error) {
	s := x.String()

	if s == "" {
		return []byte("null"), nil
	}

	return []byte(`"` + s + `"`), nil
}

func (x *HexBin) Scan(src interface{}) error {
	if src == nil {
		*x = nil
		return nil
	}

	switch s := src.(type) {
	case []byte:
		// If we got bytes, assume bytes are
		// already in hexdecimal format.
		*x = s
		return nil

	case string:
		// If we got string, it must be HEXed,
		// so we need to decode it.
		tmp, err := DecodeHexString(s)
		if err != nil {
			return err
		}
		*x = tmp
		return nil

	default:
		return errors.New("incompatible type to scan")
	}
}

func (x HexBin) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	return []byte(x), nil
}
