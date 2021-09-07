package config

import (
	"flag"
	"fmt"

	"os"

	"github.com/spf13/viper"
)

var (
	Viper    *viper.Viper
	fileName string
	Port     string
	h        bool
)

func usage() {
	fmt.Fprintf(os.Stderr, "Default Usage of ./xxx -c config.yaml\n")
	flag.PrintDefaults()
}

func initConfig() {
	flag.StringVar(&Port, "port", "5000", "监听的端口，default :5000")
	flag.StringVar(&fileName, "c", "config.yaml", "配置文件名或路径 , 默认 config.yaml")
	flag.BoolVar(&h, "h", true, "  help  ")
	flag.Usage = usage
	flag.Parse()
	if fileName != "config.yaml" {
		h = false
		flag.Usage()
	}
	Viper = viper.New()
	Viper.SetConfigFile(fileName)
	if err := Viper.ReadInConfig(); err != nil {
		os.Exit(2)
		// panic(err)
	}

}

func init() {
	initConfig()
}
