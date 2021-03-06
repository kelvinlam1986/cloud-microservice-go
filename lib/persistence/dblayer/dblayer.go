package dblayer

import (
	"cloud-microservice-go/lib/persistence"
	"cloud-microservice-go/lib/persistence/mongolayer"
)

type DBTYPE string

const (
	MONGODB DBTYPE = "mongodb"
	DOCUMENTDB DBTYPE = "documentdb"
	DYNAMODB DBTYPE = "dynamodb"
)

func NewPersistenceLayer(options DBTYPE, connection string) (persistence.DatabaseHandler, error)  {
	switch options {
	case MONGODB:
		return mongolayer.NewMongoDBLayer(connection)
	}

	return nil, nil
}