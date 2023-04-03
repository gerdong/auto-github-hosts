package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

var (
	Hosts                []string
	UpdateInterval       int
	ImmediatelyAtStartup bool
)

// Init 使用viper读取配置文件
func Init() {
	setConfigFile()
	readConfig()
	watchConfig()
}

func readConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	load()
}

func load() {
	Hosts = viper.GetStringSlice("hosts")
	UpdateInterval = viper.GetInt("update_interval")
	ImmediatelyAtStartup = viper.GetBool("immediately_at_startup")
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		load()
	})
}

func setConfigFile() {
	var file string
	flag.StringVarP(&file, "config", "c", "", "please specified config file.")
	flag.Parse()
	if file == "" {
		file = "auto-github-hosts.toml"
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("toml")
}
