package config

type RabbitMQConfiguration struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Account  string `yaml:"account"`
	Password string `yaml:"password"`
}
