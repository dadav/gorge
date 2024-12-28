package config

var (
	User                  string
	Group                 string
	ApiVersion            string
	Port                  int
	Bind                  string
	Dev                   bool
	DropPrivileges        bool
	UI                    bool
	ModulesDir            string
	ModulesScanSec        int
	Backend               string
	CORSOrigins           string
	FallbackProxyUrl      string
	NoCache               bool
	CachePrefixes         string
	CacheByFullRequestURI bool
	CacheMaxAge           int64
	ImportProxiedReleases bool
	JwtSecret             string
	TlsCertPath           string
	TlsKeyPath            string
	JwtTokenPath          string
)
