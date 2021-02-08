package config

type CassandraConfiguration struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
