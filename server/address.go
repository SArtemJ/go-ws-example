package server

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func PreparedAddressPort() (string, error) {
	var port string
	var err error

	if viper.GetString("ws.port") != "" {
		port = fmt.Sprintf("%d", viper.GetInt("ws.port"))
	} else {
		err = errors.New(fmt.Sprintf("Wrong servers address port %d", viper.GetInt("ws.port")))
	}

	return port, err
}

func PreparedAddressHost() (string, error) {
	var host string
	var err error

	if viper.GetString("ws.host") != "" {
		host = fmt.Sprintf("%s", viper.GetString("ws.host"))
	} else {
		err = errors.New(fmt.Sprintf("Wrong servers address host %d", viper.GetInt("ws.port")))
	}

	return host, err
}
