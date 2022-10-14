package config

import (
	"flag"
	"fmt"

	"os"

	"github.com/spf13/viper"
)

var (
	vp       *viper.Viper
	fileName string
	Port     string
	h        bool
	Conf     *DingtalkConfig
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
	vp = viper.New()
	vp.SetConfigFile(fileName)
	if err := vp.ReadInConfig(); err != nil {
		println(err.Error())
		os.Exit(2)
		// panic(err)
	}
	vp.Unmarshal(&Conf)

}

func init() {
	initConfig()
}
