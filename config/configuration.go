package config

type Configuration struct {
	Server ServerConfiguration `yaml:"server"`
}

func NewDefaultConfig() []byte {
	defaultConf := []byte(`
		server:
			host: "localhost"
			port: 8080
			runMode: "debug"
		`)
	return defaultConf
}
