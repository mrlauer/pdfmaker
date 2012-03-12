// package db contains generic wrappers for database access
package db

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strconv"
)

// DB is an interface for database back ends
type DB interface {
	Add(collection string, obj DBObjectWriter) error
	Update(collection string, obj DBObject) error
	Fetch(collection string, id Id, obj DBObjectWriter) error
	Delete(collection string, id Id) error
	DeleteAll(collection string) error
	Count(collection string) int
	DropDB() error
	Close()
}

// DBObject is an interface for things that can be stored in databases
type DBObject interface {
	ObjectId() Id
}

type DBObjectWriter interface {
	DBObject
	SetObjectId(id Id)
}

type Id struct {
	impl string
}

func (d Id) IsNull() bool {
	return d.impl == ""
}

func (d Id) IsValid() bool {
	return d.impl != ""
}

func (d Id) String() string {
	return d.impl
}

func MakeId(impl string) Id {
	return Id{impl}
}

// NewId generates a (presumably!) unique docId.
// It uses rand to create a type 4 uuid, then base64s it
func NewId() (Id, error) {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	impl := base64.URLEncoding.EncodeToString(b)
	return MakeId(impl), err
}

func MakeIdInt(val int) Id {
	impl := strconv.FormatInt(int64(val), 10)
	return Id{impl}
}

func (d Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Id) UnmarshalJSON(data []byte) error {
	var impl string
	err := json.Unmarshal(data, &impl)
	if err == nil {
		*d = MakeId(impl)
	}
	return err
}
