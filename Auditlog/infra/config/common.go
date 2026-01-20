package config

type SvcConf struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}
