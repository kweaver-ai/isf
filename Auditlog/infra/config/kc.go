package config

type KcConf struct {
	Private SvcConf `yaml:"private"`
	Public  SvcConf `yaml:"public"`
}
