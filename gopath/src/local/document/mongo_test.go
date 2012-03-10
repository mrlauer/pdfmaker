package document

import (
	"testing"
)

func TestMongoDB(t *testing.T) {
	mdb, err := CreateMongoDB("localhost", "testdb")

	if err != nil {
		t.Errorf("Could not create DB: %q", err.Error())
	}
	defer mdb.Close()
	defer mdb.DeleteAll()

	N := 20
	defaultText := "Hello World"
	defaultHeight := `10"`
	defaultFontSize := 12.5
	ids := []DocId{}
	for i := 1; i <= N; i++ {
		var doc Document
		doc.Text = defaultText
		doc.PageHeight, _ = LengthFromString(defaultHeight)
		doc.FontSize = LengthFromPoints(defaultFontSize)
		mdb.Add(&doc)
		for _, oldId := range ids {
			if oldId == doc.Id {
				t.Errorf("repeated id %q", doc.Id.String())
			}
		}
		ids = append(ids, doc.Id)
	}

	fetch := func(id int) (Document, error) {
		return mdb.Fetch(ids[id])
	}

	{
		doc, err := fetch(2)
		if err != nil {
			t.Errorf("Could not find document")
		}
		if doc.Text != defaultText {
			t.Errorf("Text was not correct (%q)", doc.Text)
		}
		if pts := doc.FontSize.Points(); pts != defaultFontSize {
			t.Errorf("Font size was %g\n", pts)
		}
		if ht := doc.PageHeight.String(); ht != defaultHeight {
			t.Errorf("Height was %q\n", ht)
		}
		testText := "Friends! Romans! Countrymen!"
		doc.Text = testText
		err = mdb.Update(&doc)
		if err != nil {
			t.Errorf("Could not update document")
		}
		doc2, err := fetch(2)
		if err != nil {
			t.Errorf("Could not find document")
		}
		if doc2.Text != testText {
			t.Errorf("Text was not updated (%q)", doc2.Text)
		}
	}
	{
		newId, _ := NewDocId()
		_, err := mdb.Fetch(newId)
		if err == nil {
			t.Errorf("Found nonexistent document")
		}
	}

	{
		mdb.Delete(ids[3])
		_, err = fetch(3)
		if err == nil {
			t.Errorf("Found nonexistent document")
		}
		ct := mdb.Count()
		if ct != N-1 {
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
