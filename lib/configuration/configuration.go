package configuration

import (
	"cloud-microservice-go/lib/persistence/dblayer"
	"encoding/json"
	"fmt"
	"os"
)

var (
	DBTypeDefault = dblayer.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://172.0.0.1"
	RestfulEPDefault = "localhost:8181"
)

type ServiceConfig struct {
	DatabaseType dblayer.DBTYPE `json:"databasetype"`
	DBConnection string `json:"dbconnection"`
	RestfulEndPoint string `json:"restfulapi_endpoint"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}

	err = json.NewDecoder(file).Decode(&conf)
	return conf, err
}