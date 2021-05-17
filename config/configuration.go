package config

type Configuration struct {
	Server   ServerConfiguration   `yaml:"server"`
	Rabbitmq RabbitMQConfiguration `yaml:"rabbitmq"`
	Fluentd  FluentdConfiguration  `yaml:"fluentd"`
	Consumer ConsumerConfiguration `yaml:"consumer"`
	Jaeger   JaegerConfiguration   `yaml:"jaeger"`
	Db       DBonfiguration        `yaml:"db"`
	Redis    RedisConfiguration    `yaml:"redis"`
}

func NewDefaultConfig() []byte {
	defaultConf := []byte(`
		server:
			host: "127.0.0.1"
			port: 8080
			runMode: "debug"
		rabbitmq:
			host: "127.0.0.1"
			port: 5672
			account: "guest"
			password: "guest"
		cassandra:
			host: "127.0.0.1"
			port: 9042
		fluentd:
			host: "127.0.0.1"
			port: 24224
		consumer:
			type: "ptt"
		jaeger:
			open: true
			host: "127.0.0.1"
			port: 6831
		db:
			host: "127.0.0.1"
			port: 5432
			user: "postgres"
			password: "password"
			dbname: "backend_admin"
			sslmode: "disable"
			sslcert: ""
			sslkey: ""
			sslrootcert: ""
		redis:
			host: "127.0.0.1"
			port: 6379
			password: ""
			db: 0
		`)
	return defaultConf
}
