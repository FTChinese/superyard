package reader

// AccountKind is an enumeration of login method.
type AccountKind int

// Allowed values for AccountKind
const (
	AccountKindNull AccountKind = iota
	AccountKindFtc
	AccountKindWx
	AccountKindLinked
)

var accountKindNames = [...]string{
	"",
	"ftc",
	"wechat",
	"linked",
}

var accountKindStrings = map[AccountKind]string{
	0: accountKindNames[0],
	1: accountKindNames[1],
	2: accountKindNames[2],
	3: accountKindNames[3],
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
