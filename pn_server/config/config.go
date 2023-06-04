package config

// 代码来自:https://zhwt.github.io/yaml-to-go/

type Config struct {
	RegisterPrivateKey string `yaml:"registerPrivateKey"`
	AdminContact       string `yaml:"adminContact"`
	Version            string `yaml:"version"`
	SessionExpire      string `yaml:"sessionExpire"`
	Web                Web    `yaml:"web"`
	Db                 Db     `yaml:"db"`
	Log                Log    `yaml:"log"`
}
type Web struct {
	Addr string `yaml:"addr"`
}
type Leveldb struct {
	Path string `yaml:"path"`
}
type Db struct {
	Use     string  `yaml:"use"`
	Leveldb Leveldb `yaml:"leveldb"`
}
type Log struct {
	Env        string `yaml:"env"`
	Path       string `yaml:"path"`
	Encoding   string `yaml:"encoding"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
}
