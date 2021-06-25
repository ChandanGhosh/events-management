package persistence

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DatabaseHandler interface {
	AddEvent(Event) ([]byte, error)
	FindEvent([]byte) (Event, error)
	FindEventByName(string) (Event, error)
	FindAllEventsAvailable() ([]Event, error)
}

const (
	DB     = "myevents"
	USERS  = "users"
	EVENTS = "events"
)

type MongoDBLayer struct {
	session *mgo.Session
}

func NewMongoDBLayer(connection string) (*MongoDBLayer, error) {
	s, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}
	return &MongoDBLayer{
		session: s,
	}, err
}

func (mgoLayer *MongoDBLayer) getFreshSession() *mgo.Session {
	return mgoLayer.session.Copy()
}

func (mgoLayer *MongoDBLayer) AddEvent(e Event) ([]byte, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()

	if !e.ID.Valid() {
		e.ID = bson.NewObjectId()
	}

	return []byte(e.ID), s.DB(DB).C(EVENTS).Insert(e)
}

func (mgoLayer *MongoDBLayer) FindEvent(id []byte) (Event, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	var e Event
	return e, s.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&e)
}

func (mgoLayer *MongoDBLayer) FindEventByName(name string) (Event, error) {
	s := mgoLayer.getFreshSession()
	defer s.Clone()

	var e Event
	err := s.DB(DB).C(EVENTS).Find(bson.M{"name": name}).One(&e)
	return e, err
}

func (mgoLayer *MongoDBLayer) FindAllEventsAvailable() ([]Event, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()

	var events []Event
	err := s.DB(DB).C(EVENTS).Find(nil).All(&events)
	return events, err
}
