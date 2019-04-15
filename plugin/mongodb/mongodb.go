package mongodb

import (
	"gopkg.in/mgo.v2"
)

// Reporter knows how to log errors to a MongoDB collection
type Reporter struct {
	session    *mgo.Session
	collection string
}

// New returns a fully populated mongodb.Reporter
func New(url string, collection string) *Reporter {

	result := Reporter{}

	session, err := mgo.Dial(url)

	if err != nil {
		panic("derp.MongoDbReporter: Unable to connect to database at: " + url)
	}

}
