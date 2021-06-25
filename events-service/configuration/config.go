package configuration

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chandanghosh/events-management/events-service/dblayer"
)

var (
	DBTypeDefault       = dblayer.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://127.0.0.1"
	RestfulEPDefault    = "localhost:8181"
)

type ServiceConfig struct {
	DatabaseType    dblayer.DBTYPE `json:"databasetype"`
	DBConnection    string         `json:"dbconnection"`
	RestfulEndpoint string         `json:"restfulendpoint"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DatabaseType:    DBTypeDefault,
		DBConnection:    DBConnectionDefault,
		RestfulEndpoint: RestfulEPDefault,
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}
	err = json.NewDecoder(file).Decode(&conf)
	return conf, err
}
