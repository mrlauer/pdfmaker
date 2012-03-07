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
	defaultText := "Hello World"
	defaultHeight := `10"`
	defaultFontSize := 12.5
	for i := 1; i<=N; i++ {
		var doc Document
		doc.Text = defaultText
		doc.PageHeight, _ = LengthFromString(defaultHeight)
		doc.FontSize = LengthFromPoints(defaultFontSize)
		mdb.Add(&doc)
		if doc.Id != DocId(i) {
			t.Errorf("Document %d has id %d", i, doc.Id)
		}
	}

	{
		doc, err := mdb.Fetch(2)
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
