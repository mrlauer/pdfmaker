package document

import (
	"errors"
	"fmt"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
)

// Getter and setter for Document Length
// MarshalJSON uses the defining string.
func (l Length) GetBSON() (interface{}, error) {
	return l.String(), nil
}

// UnmarshalJSON uses the defining string.
func (l *Length) SetBSON(raw bson.Raw) error {
	var def string
	err := raw.Unmarshal(&def)
	if err == nil {
		*l, err = LengthFromString(def)
	}
	return err
}

func (d DocId) GetBSON() (interface{}, error) {
	return d.impl, nil
}

// UnmarshalJSON uses the defining string.
func (d *DocId) SetBSON(raw bson.Raw) error {
	var def string
	err := raw.Unmarshal(&def)
	if err == nil {
		*d = MakeDocId(def)
	}
	return err
}

func bogusCheck() {
	var l Length
	var _ bson.Getter = l
	var _ bson.Setter = &l
}

type MongoDBDoc struct {
	Id  bson.ObjectId `bson:"_id,omitempty"`
	Doc Document
}

type MongoDB struct {
	Session *mgo.Session
	DBName  string
}

func CreateMongoDB(host string, dbname string) (*MongoDB, error) {
	session, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}
	m := &MongoDB{Session: session, DBName: dbname}
	// Initialize the index for ids
	idx := mgo.Index{}
	idx.Key = []string{"-doc.id"}
	idx.Unique = true
	m.Documents().EnsureIndex(idx)
	return m, nil
}

func (m *MongoDB) Close() {
	m.Session.Close()
}

func (m *MongoDB) Database() *mgo.Database {
	return m.Session.DB(m.DBName)
}

func (m *MongoDB) Documents() *mgo.Collection {
	return m.Database().C("documents")
}

func (m *MongoDB) Count() int {
	c := m.Documents()
	n, _ := c.Find(bson.M{}).Count()
	return n
}


func (m *MongoDB) Add(doc *Document) {
	// Use optimistic loop strategy. Maybe use a counter instead at some point.
	c := m.Documents()
	toadd := MongoDBDoc{Doc: *doc}
	for true {
		// MapReduce!
		job := mgo.MapReduce {
			Map: `function() { 
				var i = parseInt(this.doc.id); 
				if(!isNaN(i)) {
					emit(0, i)
				}
			}`,
			Reduce: `function(key, values) { return Math.max.apply(null, values); }`,
		}
		var result []struct{ Id int "_id"; Value int }
		_, err := c.Find(nil).MapReduce(job, &result)
		if err != nil {
			fmt.Printf("MapReduce error: %q\n", err.Error())
		}
		id := 1
		if err == nil && len(result) > 0 {
			id = result[0].Value + 1
		}
		toadd.Doc.Id = MakeDocIdInt(id)
		err = c.Insert(&toadd)
		doc.Id = toadd.Doc.Id
		if err == nil {
			return
		}
	}
}

func (m *MongoDB) Update(doc *Document) error {
	c := m.Documents()
	mdoc := MongoDBDoc{Doc: *doc}
	err := c.Update(bson.M{"doc.id": doc.Id}, &mdoc)
	// Proper checking?
	if err != nil {
		return errors.New("document does not exist")
	}
	return nil
}

func (m *MongoDB) Fetch(id DocId) (Document, error) {
	c := m.Documents()
	mdoc := MongoDBDoc{}
	err := c.Find(bson.M{"doc.id": id}).One(&mdoc)
	if err != nil {
		return mdoc.Doc, err
	}
	return mdoc.Doc, nil
}

func (m *MongoDB) Delete(id DocId) error {
	c := m.Documents()
	err := c.Remove(bson.M{"doc.id": id})
	return err
}

func (m *MongoDB) DeleteAll() error {
	c := m.Documents()
	err := c.RemoveAll(bson.M{})
	return err
}

func (m *MongoDB) DropDB() error {
	err := m.Database().DropDatabase()
	return err
}
