package main

import (
	"cloud-microservice-go/lib/configuration"
	"cloud-microservice-go/lib/persistence/dblayer"
	"cloud-microservice-go/rest"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main()  {
	configPath := flag.String("conf", `.\configuration\config.json`,"flag to set the path to the configuration json file")
	flag.Parse()
	// extract configuration
	config, _ := configuration.ExtractConfiguration(*configPath)
	fmt.Println("Connecting to database...")
	dbHandler, _ := dblayer.NewPersistenceLayer(config.DatabaseType, config.DBConnection)
	logrus.Println(rest.ServeAPI(config.RestfulEndPoint, dbHandler))
}
