package mongo

import (
	"github.com/kirthiprakash/tech-talk-tekion-201804/app"
	"gopkg.in/mgo.v2"
)

var _db *mgo.Session

func init() {
	session, dialErr := mgo.Dial(app.MongoHost)
	if dialErr != nil {
		panic(dialErr)
	}
	_db = session
}

func GetConnection() *mgo.Session {
	return _db
}

func CloseConnection() {
	_db.Close()
}
