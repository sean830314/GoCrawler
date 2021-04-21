package config

type JaegerConfiguration struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Open bool   `yaml:"open"`
}
