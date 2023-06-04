package config

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// SugaredConfig 将配置文件的参数解析,比如解析时间为 time.Ticker
type SugaredConfig struct {
	*Config
	SessionExpireDuation time.Duration
}

var Global *SugaredConfig

func Init(filePath string) *SugaredConfig {
	// 初始化配置文件
	pflag.StringP("config", "c", filePath, "config file")
	pflag.Parse()
	viper.SetConfigType("yaml")
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(err)
	}
	conf := viper.GetString("config")

	fmt.Println("=====================================")
	fmt.Printf("[config file]: path %s\n", conf)
	fmt.Println("=====================================")

	viper.SetConfigFile(conf)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("load config %s fail: %v", conf, err))
	}

	// 解析初始配置
	baseConf := &Config{}
	if err := viper.Unmarshal(baseConf); err != nil {
		if err != nil {
			panic(err)
		}
	}

	// 解析 session 过期时间
	if baseConf.SessionExpire == "" {
		baseConf.SessionExpire = "1h"
	}
	d, err := time.ParseDuration(baseConf.SessionExpire)
	if err != nil {
		panic(err)
	}

	// 构造 SugaredConfig
	C := &SugaredConfig{
		Config:               baseConf,
		SessionExpireDuation: d,
	}

	Global = C

	return Global
}
