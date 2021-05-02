package consts

const (
	EnvVarPrefix         = "GO_CRAWLER"
	AllowEmptyEnv        = true
	DefaultConfigPath    = "/etc/GoCrawler/config.yaml"
	DefaultLogOutputPath = "/var/log/GoCrawler/"

	DefDBHost        = "localhost"
	DefDBPort        = "5432"
	DefDBUser        = "postgres"
	DefDBPass        = "password"
	DefDBName        = "backend_admin"
	DefDBSSLMode     = "disable"
	DefDBSSLCert     = ""
	DefDBSSLKey      = ""
	DefDBSSLRootCert = ""
)
