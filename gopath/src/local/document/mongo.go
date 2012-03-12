package document

import (
	"db"
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

func (doc Document) ObjectId() db.Id {
	return doc.Id
}

func (doc *Document) SetObjectId(id db.Id) {
	doc.Id = id
}

func bogusCheck() {
	var l Length
	var _ bson.Getter = l
	var _ bson.Setter = &l

	var doc Document
	var _ db.DBObject = doc
	var _ db.DBObjectWriter = &doc
}

type MongoDB struct {
	Database db.DB
}

func CreateMongoDB(host string, dbname string) (*MongoDB, error) {
	database, err := db.CreateMongoDB(host, dbname)
	if err != nil {
		return nil, err
	}
	return &MongoDB{database}, nil
}

func (m *MongoDB) Close() {
	m.Database.Close()
}

const docCollection = "documents"

func (m *MongoDB) Count() int {
	return m.Database.Count(docCollection)
}

func (m *MongoDB) Add(doc *Document) error {
	return m.Database.Add(docCollection, doc)
}

func (m *MongoDB) Update(doc *Document) error {
	return m.Database.Update(docCollection, doc)
}

func (m *MongoDB) Fetch(id db.Id) (Document, error) {
	var doc Document
	err := m.Database.Fetch(docCollection, id, &doc)
	return doc, err
}

func (m *MongoDB) Delete(id db.Id) error {
	return m.Database.Delete(docCollection, id)
}

func (m *MongoDB) DeleteAll() error {
	return m.Database.DeleteAll(docCollection)
}

func (m *MongoDB) DropDB() error {
	return m.Database.DropDB()
}
