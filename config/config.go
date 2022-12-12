package config

import (
	"flag"

	"fmt"
	"github.com/mlonV/tools/loger"
	"github.com/spf13/viper"
	"os"
)

var (
	vp       *viper.Viper
	fileName string
	Port     string
	h        bool
	Conf     *DingtalkConfig
	Log      *loger.Loger
)

func usage() {
	fmt.Fprintf(os.Stderr, "Default Usage of ./xxx -c config.yaml\n")
	flag.PrintDefaults()
}

func initConfig() {
	flag.StringVar(&Port, "port", "5000", "监听的端口，default :5000")
	flag.StringVar(&fileName, "c", "config.yaml", "配置文件名或路径 , 默认 config.yaml")
	flag.BoolVar(&h, "h", false, "  help  ")
	flag.Usage = usage
	flag.Parse()
	// 使LoadConfig() 可以多次执行

	if err := LaodConfig(); err != nil {
		panic(err)
	}
}

func initLog() {
	// 初始化全局打印日志，获取config.yaml配置文件里面的设置引用到Log
	Log = loger.NewLoger(
		&loger.Loger{
			ToFile:          Conf.LogSet.ToFile,
			WithFuncAndFile: Conf.LogSet.WithFuncAndFile,
			Level:           Conf.LogSet.Level,
			FileLoger: loger.FileLoger{
				FileName:    Conf.LogSet.FileName,
				FilePath:    Conf.LogSet.FilePath,
				FileMaxSize: Conf.LogSet.FileSize,
				FileSaveNum: Conf.LogSet.FileSaveNum,
			},
		},
	)
}

func LaodConfig() error {
	if h {
		flag.Usage()
		os.Exit(0)
	}
	vp = viper.New()
	vp.SetConfigFile(fileName)
	if err := vp.ReadInConfig(); err != nil {
		return err
	}
	if err := vp.Unmarshal(&Conf); err != nil {
		return err
	}
	return nil
}

func init() {
	initConfig()
	initLog()
}
