package config

type DBonfiguration struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Pass        string `yaml:"password"`
	Name        string `yaml:"dbname"`
	SSLMode     string `yaml:"sslmode"`
	SSLCert     string `yaml:"sslcert"`
	SSLKey      string `yaml:"sslkey"`
	SSLRootCert string `yaml:"sslrootcert"`
}
