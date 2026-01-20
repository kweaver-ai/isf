package config

type OssGatewayConf struct {
	Private SvcConf `yaml:"private"`
	Public  SvcConf `yaml:"public"`
}
