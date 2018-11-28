package main

import (
	"cloud-microservice-go/lib/configuration"
	"cloud-microservice-go/lib/persistence/dblayer"
	"cloud-microservice-go/eventservice/rest"
	"flag"
	"fmt"
	"log"
)

func main()  {
	configPath := flag.String("conf", `.\configuration\config.json`,"flag to set the path to the configuration json file")
	flag.Parse()
	// extract configuration
	config, _ := configuration.ExtractConfiguration(*configPath)
	fmt.Println("Connecting to database...")
	dbHandler, _ := dblayer.NewPersistenceLayer(config.DatabaseType, config.DBConnection)
	httpError, httptlsError := rest.ServeAPI(config.RestfulEndPoint, config.RestfulTLSEndPint, dbHandler)
	fmt.Println("serve event service at http port: " + config.RestfulEndPoint)

	select {
	case err := <- httpError:
		log.Fatal("HTTP Error:", err)
	case err := <- httptlsError:
		log.Fatal("HTTPS Error:", err)
	}
}
