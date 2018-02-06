package configuration

import (
	"encoding/json"
	"fmt"
	"os"

	dblayer "github.com/omgitsotis/user-service/dblayer"
)

var (
	DBTypeDefault      = dblayer.MOCKDB
	RestfulEPDefault   = "localhost:8080"
	DefaultAMQPBrooker = "test"
)

type ServiceConfig struct {
	DatabaseLayer dblayer.DBType `json:"database_type"`
	RestfulEP     string         `json:"endpoint"`
	AMQPBrooker   string         `json:"amqp_brooker"`
}

func GetConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{DBTypeDefault, RestfulEPDefault, DefaultAMQPBrooker}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found, using defaults")
		return conf, err
	}

	err = json.NewDecoder(file).Decode(&conf)
	return conf, err
}
