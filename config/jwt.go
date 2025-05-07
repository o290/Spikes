package config

type Jwt struct {
	SignKey    string
	ExpireTime int
	Issuer     string
}
