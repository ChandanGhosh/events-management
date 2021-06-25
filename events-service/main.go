package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/chandanghosh/events-management/events-service/configuration"
	"github.com/chandanghosh/events-management/events-service/dblayer"
	"github.com/chandanghosh/events-management/events-service/handler"
	"github.com/chandanghosh/events-management/events-service/persistence"
	"github.com/gorilla/mux"
)

const (
	DB_CONN_STR = "mongo://localhost:27017/admin"
	APP_PORT    = "8181"
)

func ServeAPI(endpoint string, dbhandler persistence.DatabaseHandler) error {

	handler := handler.NewEventServiceHandler(dbhandler)

	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("events").Subrouter()

	eventsrouter.Methods("GET").Path("/{searchcriteria}/{search}").HandlerFunc(handler.FindEventsHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.AllEventsHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.NewEventHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = APP_PORT
	}
	return http.ListenAndServe(endpoint, r)
}

func main() {
	confPath := flag.String("conf", `./configuration/config.json`, "flag to set the path to the configuration file.")
	flag.Parse()
	conf, err := configuration.ExtractConfiguration(*confPath)
	if err != nil {
		log.Fatalf("The configuration can not be extracted. error: %s", err)
		os.Exit(1)
	}
	dbhandler, _ := dblayer.NewPersistenceLayer(conf.DatabaseType, conf.DBConnection)
	log.Fatalln(ServeAPI(conf.RestfulEndpoint, dbhandler))
}
