package dblayer

import (
	"github.com/chandanghosh/events-management/events-service/persistence"
)

type DBTYPE string

const (
	MONGODB  = "mongodb"
	DYNAMODB = "dynamodb"
)

func NewPersistenceLayer(option DBTYPE, connstr string) (*persistence.MongoDBLayer, error) {
	switch option {
	case MONGODB:
		return persistence.NewMongoDBLayer(connstr)
	}
	return nil, nil
}
