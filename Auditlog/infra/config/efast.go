package config

type EFastConf struct {
	Private SvcConf `yaml:"private"`
	Public  SvcConf `yaml:"public"`
}
