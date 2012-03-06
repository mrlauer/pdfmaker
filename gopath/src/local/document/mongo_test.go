package document

import(
	"testing"
)


func TestMongoDB(t *testing.T) {
	mdb, err := CreateMongoDB("localhost", "testdb")

	if err != nil {
		t.Errorf("Could not create DB: %q", err.Error())
	}
	defer mdb.Close()
	defer mdb.DeleteAll()

	N := 4
	for i := 1; i<=N; i++ {
		var doc Document
		mdb.Add(&doc)
		if doc.Id != i {
			t.Errorf("Document %d has id %d", i, doc.Id)
		}
	}

	{
		doc, err := mdb.Fetch(2)
		if err != nil {
			t.Errorf("Could not find document")
		}
		testText := "Friends! Romans! Countrymen!"
		doc.Text = testText
		err = mdb.Update(&doc)
		if err != nil {
			t.Errorf("Could not update document")
		}
		doc2, err := mdb.Fetch(2)
		if err != nil {
			t.Errorf("Could not find document")
		}
		if doc2.Text != testText {
			t.Errorf("Text was not updated (%q)", doc2.Text)
		}
	}
	{
		_, err := mdb.Fetch(47)
		if err == nil {
			t.Errorf("Found nonexistent document")
		}
	}

	{
		mdb.Delete(3)
		_, err = mdb.Fetch(3)
		if err == nil {
			t.Errorf("Found nonexistent document")
		}
		ct := mdb.Count()
		if ct != 3 {
			t.Errorf("Wrong count after remove: %d", ct)
		}
	}
	{
		mdb.DeleteAll()
		ct := mdb.Count()
		if ct != 0 {
			t.Errorf("Wrong count after remove all: %d", ct)
		}
	}

}
