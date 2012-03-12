package db

import (
	"errors"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
)

type MongoDB struct {
	Session *mgo.Session
	DBName  string
}

type ObjectHolder struct {
	Object interface{}
}

type MongoDBObject struct {
	Id     Id `bson:"_id,omitempty"`
	Object interface{}
}

func CreateMongoDB(host string, dbname string) (*MongoDB, error) {
	session, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}
	m := &MongoDB{Session: session, DBName: dbname}
	return m, nil
}

func (m *MongoDB) Close() {
	m.Session.Close()
}

func (m *MongoDB) Database() *mgo.Database {
	return m.Session.DB(m.DBName)
}

func (m *MongoDB) Collection(collection string) *mgo.Collection {
	return m.Database().C(collection)
}

func (m *MongoDB) Count(collection string) int {
	c := m.Collection(collection)
	n, _ := c.Find(bson.M{}).Count()
	return n
}

func (m *MongoDB) Add(collection string, obj DBObjectWriter) error {
	c := m.Collection(collection)
	toadd := MongoDBObject{Object: obj}
	for true {
		id, err := NewId()
		if err != nil {
			return errors.New("Could not create new Id: " + err.Error())
		}
		toadd.Id = id
		obj.SetObjectId(id)
		err = c.Insert(&toadd)
		if err != nil {
			// Should check that it was a duplicate error
			continue
		}
		break
	}
	return nil
}

func (m *MongoDB) Update(collection string, obj DBObject) error {
	c := m.Collection(collection)
	mdoc := MongoDBObject{Id: obj.ObjectId(), Object: obj}
	err := c.Update(bson.M{"_id": obj.ObjectId()}, &mdoc)
	// Proper checking?
	if err != nil {
		return errors.New("document does not exist")
	}
	return nil
}

func (m *MongoDB) Fetch(collection string, id Id, obj DBObjectWriter) error {
	c := m.Collection(collection)
	mdoc := MongoDBObject{Id: id}
	err := c.Find(bson.M{"_id": id}).One(&mdoc)
	if err != nil {
		return err
	}
	// The value will now be in bson.M form in the Object field. Reserialize it
	// to the receiver. It would be more clever to have a wrapper in mdoc around
	// the interface that does the serialization.
	bytes, err := bson.Marshal(mdoc.Object)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bytes, obj)
	if err != nil {
		return err
	}
	return err
}

func (m *MongoDB) Delete(collection string, id Id) error {
	c := m.Collection(collection)
	err := c.Remove(bson.M{"_id": id})
	return err
}

func (m *MongoDB) DeleteAll(collection string) error {
	c := m.Collection(collection)
	err := c.RemoveAll(bson.M{})
	return err
}

func (m *MongoDB) DropDB() error {
	err := m.Database().DropDatabase()
	return err
}

func (d Id) GetBSON() (interface{}, error) {
	return d.impl, nil
}

// UnmarshalJSON uses the defining string.
func (d *Id) SetBSON(raw bson.Raw) error {
	var def string
	err := raw.Unmarshal(&def)
	if err == nil {
		*d = MakeId(def)
	}
	return err
}

func bogusCheckerFunction() {
	var db MongoDB
	var _ DB = &db
	var id Id
	var _ bson.Setter = &id
	var _ bson.Getter = id
}
