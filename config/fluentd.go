package config

type FluentdConfiguration struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
