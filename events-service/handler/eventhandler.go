package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chandanghosh/events-management/events-service/persistence"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type EventServiceHandler struct {
	dbhandler persistence.DatabaseHandler
}

func NewEventServiceHandler(databaseHandler persistence.DatabaseHandler) *EventServiceHandler {
	return &EventServiceHandler{
		dbhandler: databaseHandler,
	}
}

// FindEventsHandler handles searchcriteria and searchterm.
// GET /events/{SearchCriteria}/{search}
func (eh *EventServiceHandler) FindEventsHandler(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	searchCriteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{error: No search criteria found, you can either search by id or by name}`)
		return
	}
	searchterm, ok := vars["search"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{error: No search term found. You can search by some id or name of the event}`)
		return
	}

	var event persistence.Event
	var err error
	var id []byte

	switch strings.ToLower(searchCriteria) {
	case "name":
		event, err = eh.dbhandler.FindEventByName(searchterm)
	case "id":
		if id, err = hex.DecodeString(searchterm); err == nil {
			event, err = eh.dbhandler.FindEvent(id)
		}
	}

	if err != nil {
		fmt.Fprintf(w, "{error: %s}", err)
		return
	}
	w.Header().Set("content-type", "application/json;chartset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh *EventServiceHandler) AllEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbhandler.FindAllEventsAvailable()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: %s}", err)
		return
	}
	w.Header().Set("content-type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&events)
}

func (eh *EventServiceHandler) NewEventHandler(w http.ResponseWriter, r *http.Request) {
	var event persistence.Event
	var err error
	json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{error: Error occured while decoding the event %s}`, err)
		return
	}

	id, err := eh.dbhandler.AddEvent(event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{error: error occure saving event %s }`, err)
		return
	}
	w.WriteHeader(201)
	event.ID = bson.ObjectId(id)
	json.NewEncoder(w).Encode(&event)
}
