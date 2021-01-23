package config

type Configuration struct {
	Server ServerConfiguration `yaml:"server"`
}

func NewDefaultConfig() []byte {
	defaultConf := []byte(`
		server:
			host: "127.0.0.1"
			port: 8080
		`)
	return defaultConf
}

// func InitConfiguration() {
// 	viper.SetConfigName("config")
// 	viper.AddConfigPath(".")
// 	var configuration config.Configuration

// 	if err := viper.ReadInConfig(); err != nil {
// 		log.Fatalf("Error reading config file, %s", err)
// 	}
// 	err := viper.Unmarshal(&configuration)
// 	if err != nil {
// 		log.Fatalf("unable to decode into struct, %v", err)
// 	}
// }