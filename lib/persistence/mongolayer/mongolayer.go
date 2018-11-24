package mongolayer

import (
	"cloud-microservice-go/lib/persistence"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB = "myevents"
	USERS = "users"
	EVENTS = "events"
)

type MongoDBLayer struct {
	session *mgo.Session
}

func NewMongoDBLayer(connection string) (persistence.DatabaseHandler, error) {
	s, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}
	return &MongoDBLayer{
		session:s,
	}, err
}

func (mgoLayer *MongoDBLayer) AddEvent(event persistence.Event) ([]byte, error)  {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	if !event.ID.Valid() {
		event.ID = bson.NewObjectId()
	}

	if !event.Location.ID.Valid() {
		event.Location.ID = bson.NewObjectId()
	}

	return []byte(event.ID), s.DB(DB).C(EVENTS).Insert(event)
}

func (mgoLayer *MongoDBLayer) FindEvent(id []byte) (persistence.Event, error)  {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	e := persistence.Event{}
	err := s.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&e)
	return e, err
}

func (mgoLayer *MongoDBLayer) FindEventByName(name string) (persistence.Event, error)  {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	e := persistence.Event{}
	err := s.DB(DB).C(EVENTS).Find(bson.M{"name": name}).One(&e)
	return e, err
}

func (mgoLayer *MongoDBLayer) FindAllAvailableEvents()  ([]persistence.Event, error)  {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	events := []persistence.Event{}
	err := s.DB(DB).C(EVENTS).Find(nil).All(&events)
	return events, err
}

func (mgoLayer *MongoDBLayer) getFreshSession() *mgo.Session {
	return mgoLayer.session.Copy()
}