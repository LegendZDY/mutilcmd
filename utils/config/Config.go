package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `mapstructure:"app" json:"app" yaml:"app"`
	Log    LogConfig    `mapstructure:"log" json:"log" yaml:"log"`
}

// 读取Yaml配置文件，并转换成Config对象  struct结构
func (config *Config) InitializeConfig() *Config {
	//获取项目的执行路径
	homeDir, _err := os.UserHomeDir()
	if _err != nil {
		panic(_err)
	}

	configPath := homeDir + "/.mutilcmd/config"
	configName := "config"
	ConfigType := "yaml"

	vip := viper.New()
	vip.AddConfigPath(configPath)           //设置读取的文件路径
	vip.SetConfigName(configName)           //设置读取的文件名
	vip.SetConfigType(ConfigType)           //设置文件的类型
	//尝试进行配置读取
	if err := vip.ReadInConfig(); err != nil {
		panic(err)
	}
	// 监听配置文件
	vip.WatchConfig()
	vip.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		// 重载配置
		if err := vip.Unmarshal(&config); err != nil {
			fmt.Println(err)
		}
	})

	err := vip.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	return config
}

