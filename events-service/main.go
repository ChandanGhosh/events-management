package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/chandanghosh/events-management/events-service/configuration"
	"github.com/chandanghosh/events-management/events-service/dblayer"
	"github.com/chandanghosh/events-management/events-service/handler"
	"github.com/chandanghosh/events-management/events-service/persistence"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"

	"github.com/chandanghosh/events-management/contracts/lib/msgqueue"
	"github.com/chandanghosh/events-management/contracts/lib/msgqueue/mqp"
)

func ServeAPI(httpEndpoint, httpsEndpoint string, dbhandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) (chan error, chan error) {

	handler := handler.NewEventServiceHandler(dbhandler, eventEmitter)

	r := mux.NewRouter()
	eventsRouter := r.PathPrefix("events").Subrouter()

	eventsRouter.Methods("GET").Path("/{searchcriteria}/{search}").HandlerFunc(handler.FindEventsHandler)
	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.AllEventsHandler)
	eventsRouter.Methods("POST").Path("").HandlerFunc(handler.NewEventHandler)

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
	// amqp_url := os.Getenv("AMQP_URL")
	// if amqp_url == "" {
	// 	amqp_url = "amqp://guest:guest@localhost:5672"
	// }
	// conn, err := amqp.Dial(amqp_url)
	// if err != nil {
	// 	log.Println("Error connecting to broker " + err.Error())
	// }
	// defer conn.Close()
	// channel, err := conn.Channel()
	// if err != nil {
	// 	panic("Could not open a channel on the broker " + err.Error())
	// }

	confPath := flag.String("conf", `./configuration/config.json`, "flag to set the path to the configuration file.")
	flag.Parse()
	conf, err := configuration.ExtractConfiguration(*confPath)
	if err != nil {
		log.Fatalf("The configuration can not be extracted. error: %s", err)
	}
	conn, err := amqp.Dial(conf.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}

	emitter, err := mqp.NewAMQPEventEmitter(conn)
	if err != nil {
		panic(err)
	}

	dbhandler, _ := dblayer.NewPersistenceLayer(conf.DatabaseType, conf.DBConnection)
	httpErr, httpsErr := ServeAPI(conf.RestfulEndpoint, conf.RestfulTLSEndpoint, dbhandler, emitter)

	select {
	case err = <-httpErr:
		log.Fatal("HTTP error: ", err)
	case err = <-httpsErr:
		log.Fatal("HTTPS error: ", err)
	}
}
