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

func ServeAPI(httpEndpoint, httpsEndpoint string, dbhandler persistence.DatabaseHandler) (chan error, chan error) {

	handler := handler.NewEventServiceHandler(dbhandler)

	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("events").Subrouter()

	eventsrouter.Methods("GET").Path("/{searchcriteria}/{search}").HandlerFunc(handler.FindEventsHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.AllEventsHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.NewEventHandler)

	httpError := make(chan error)
	httpsError := make(chan error)
	go func() {
		httpError <- http.ListenAndServe(httpEndpoint, r)
	}()
	go func() {
		httpsError <- http.ListenAndServeTLS(httpsEndpoint, "cert.pem", "key.pem", r)
	}()
	return httpError, httpsError
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
	httpErr, httpsErr := ServeAPI(conf.RestfulEndpoint, conf.RestfulTLSEndpoint, dbhandler)

	select {
	case err = <-httpErr:
		log.Fatal("HTTP error: ", err)
	case err = <-httpsErr:
		log.Fatal("HTTPS error: ", err)
	}
}
