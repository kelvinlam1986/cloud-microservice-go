package configuration

import (
	"cloud-microservice-go/lib/persistence/dblayer"
	"encoding/json"
	"fmt"
	"os"
)

var (
	DBTypeDefault = dblayer.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://localhost:27017"
	RestfulEPDefault = "localhost:8282"
	RestfulTLSEPDefault = "localhost:9191"
	MessageBrokerTypeDefault = "amqp"
	AMQPMessageBrokerDefault = "amqp://guest:guest@localhost:5672"
)

type ServiceConfig struct {
	DatabaseType dblayer.DBTYPE `json:"databasetype"`
	DBConnection string `json:"dbconnection"`
	RestfulEndPoint string `json:"restfulapi_endpoint"`
	RestfulTLSEndPint string `json:"restfulapi-tlsendpoint"`
	MessageBrokerType string `json:"message_broker_type"`
	AMQPMessageBroker string `json:"amqp_message_broker"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		RestfulTLSEPDefault,
		MessageBrokerTypeDefault,
		AMQPMessageBrokerDefault,
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}

	err = json.NewDecoder(file).Decode(&conf)

	if v := os.Getenv("AMQP_BROKER_URL"); v != "" {
		conf.MessageBrokerType = "amqp"
		conf.AMQPMessageBroker = v
	}

	return conf, err
}