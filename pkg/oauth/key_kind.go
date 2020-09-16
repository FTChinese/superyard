package oauth

// KeyUsage tells the kind of an access token
type KeyKind string

const (
	KeyKindApp      KeyKind = "app"      // Used by an app.
	KeyKindPersonal KeyKind = "personal" // Used by human.
)
