package main

import (
	"fmt"

	"github.com/farmako/api"
	"github.com/farmako/config"
)

func main() {
	config, err := config.SetupEnv()
	if err != nil {
		fmt.Errorf(err.Error())
	}

	api.StartServer(config)
}
