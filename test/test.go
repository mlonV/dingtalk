package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type AlertNotifyList struct {
	AlertJobList `mapstructure:"dingtalk"`
}
type AlertJobList struct {
	At     string `mapstructure:"at"`
	Job    string `mapstructure:"job"`
	Secret string `mapstructure:"secret"`
	Url    string `mapstructure:"url"`
}

func main() {

	// var Viper *viper.Viper
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	viper.SetConfigName("dingtalk.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("err := ", err)
	}
	// toml := viper.Get("groups")
	anl := &AlertJobList{}
	err := viper.Unmarshal(&AlertJobList{})
	// toml := viper.Get("dingtalk")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(anl)

	// switch toml.(type) {
	// case map[string]map[string]string:
	// 	// fmt.Println(toml)
	// 	for k, v := range toml {
	// 		fmt.Println(k, v)
	// 	}
	// default:
	// 	fmt.Println("default")
	// }

}
