package config

type DocCenterConf struct {
	Private SvcConf `yaml:"private"`
	Public  SvcConf `yaml:"public"`
}
