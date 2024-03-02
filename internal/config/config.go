package config

var (
	ApiVersion            string
	Port                  int
	Bind                  string
	Dev                   bool
	ModulesDir            string
	Backend               string
	CORSOrigins           string
	FallbackProxyUrl      string
	NoCache               bool
	CachePrefixes         string
	CacheDir              string
	CacheMaxAge           int64
	ImportProxiedReleases bool
	JwtSecret             string
	TlsCertPath           string
	TlsKeyPath            string
	JwtTokenPath          string
)
