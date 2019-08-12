package util

type BuildConfig struct {
	IsProduction bool   `json:"isProduction"`
	Version      string `json:"version"`
	BuiltAt      string `json:"builtAt"`
}
