package config

type ServerConfiguration struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	RunMode string `yaml:"runMode"`
}
