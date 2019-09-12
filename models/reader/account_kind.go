package reader

// AccountKind is an enumeration of login method.
type AccountKind int

// Allowed values for AccountKind
const (
	AccountKindFtc AccountKind = iota
	AccountKindWx
)

var accountKindNames = [...]string{
	"ftc",
	"wechat",
}

var accountKindStrings = map[AccountKind]string{
	0: accountKindNames[0],
	1: accountKindNames[1],
}

func (x AccountKind) String() string {
	if str, ok := accountKindStrings[x]; ok {
		return str
	}

	return ""
}

// MarshalJSON implements the Marshaler interface
func (x AccountKind) MarshalJSON() ([]byte, error) {
	s := x.String()

	if s == "" {
		return []byte("null"), nil
	}

	return []byte(`"` + s + `"`), nil
}
