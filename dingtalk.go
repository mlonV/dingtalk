package main

import (
	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/route"
)

// "github.com/mlonV/dingtalk/config"

func main() {

	r := route.RegisterRoutes()
	r.Run(":" + config.Port)

}
