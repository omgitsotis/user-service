package main

import (
    "log"
    "flag"


    configuration "github.com/omgitsotis/user-service/configuration"
    dblayer "github.com/omgitsotis/user-service/dblayer"
    client "github.com/omgitsotis/user-service/client"
)

func main() {
    confPath := flag.String("conf", `configuration\config.json`, "floag to set the path of the configuration json file")
    flag.Parse()
    config, _ := configuration.GetConfiguration(*confPath)
    dbHandler, _ := dblayer.NewPersistenceLayer(config.DatabaseLayer, "")

    log.Fatal(client.ServeAPI(dbHandler, config.RestfulEP))
}
