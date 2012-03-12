package db

import (
	"testing"
)

type Thing struct {
	Id	 Id
	Data string
}

func (t *Thing)ObjectId() Id {
	return t.Id
}

func (t *Thing)SetObjectId(id Id) {
	t.Id = id
}

func TestMongoDB(t *testing.T) {
	var db DB
	var err error
	db, err = CreateMongoDB("localhost", "testdb")
	if err != nil {
		t.Errorf("Error creating DB: %q", err.Error())
	}
	defer db.Close()
	defer db.DropDB()

	thing := Thing{Data:"Initial"}
	coll := "things"
	err = db.Add(coll, &thing)
	if err != nil {
		t.Errorf("Error adding: %q", err.Error())
	}
	id := thing.Id
	if id.IsNull() {
		t.Errorf("Null Id in added object")
	}

	var thing2 Thing
	err = db.Fetch(coll, id, &thing2)
	if err != nil {
		t.Errorf("Could not fetch object: %q", err.Error())
	}
	if thing2.Data != thing.Data {
		t.Errorf("Bad data: %s", thing2.Data)
	}

	thing.Data = "Modified"
	db.Update(coll, &thing)
	var thing3 Thing
	db.Fetch(coll, id, &thing3)
	if thing3.Id != id {
		t.Errorf("Fetched id is incorrect")
	}
	if thing3.Data != "Modified" {
		t.Errorf("Modifed/fetched data incorrect: %s", thing3.Data)
	}

	c1 := db.Count(coll)
	if c1 != 1 {
		t.Errorf("Bad count %d", c1)
	}
	db.Delete(coll, id)
	c2 := db.Count(coll)
	if c2 != 0 {
		t.Errorf("Bad count %d", c2)
	}

	var thing4 Thing
	err = db.Fetch(coll, id, &thing)
	if err == nil {
		t.Errorf("object was not deleted")
	}
	if thing4.Id.IsValid() {
		t.Errorf("fetched id %q was improperly valid", thing4.Id)
	}

}
