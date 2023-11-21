package config

type DbConf struct {
	Type, Host, User, Password, Name, Port string
	MaxIdle, MaxOpen                       int
	SingularTable                          bool
	PreferSimpleProtocol                   bool
}
type GinConfig struct {
	Mode        string
	Port        string
	Domain      string
	FrontDomain string `yaml:"front_domain"`
	FrontDir    string `yaml:"front_dir"`
}
type RpcConfig struct {
	Host string
	Mode string
}

type LogConfig struct {
	Mode       string
	Host       string
	Db         string
	Collection string
	User       string
	Password   string
}
type CConfig struct {
	Host     string
	Password string
	Db       int
}
type QnConfig struct {
	Host string
	AK   string
	SK   string
	SC   string
	ST   string
}
type QueueConfig struct {
	Producer string
	Customer []string
}

type AppConf struct {
	Owner   string
	BackRpc []string `yaml:"back_rpc"`
	Rpc     *RpcConfig
	Gin     *GinConfig
	Db      *DbConf
	Log     *LogConfig
	Cache   *CConfig
	Queue   *QueueConfig
	Qn      *QnConfig
	Custom  map[string]any
}
