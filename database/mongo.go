package database

import (
	"fmt"

	"github.com/jamesonwilliams/golang-sso-google/auth"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const databaseHostname = "localhost"
const databaseName = "sso"
const usersCollectionKey = "user"

// MongoDBConnection Encapsulates a connection to a database.
type MongoDBConnection struct {
	session *mgo.Session
}

// SaveUser register a user so we know that we saw that user already.
func (mongoDb MongoDBConnection) SaveUser(user *auth.User) error {
	mongoDb.session = mongoDb.GetSession()
	defer mongoDb.session.Close()
	if _, err := mongoDb.LoadUser(user.Email); err == nil {
		return fmt.Errorf("User already exists!")
	}
	users := mongoDb.session.DB(databaseName).C(usersCollectionKey)
	err := users.Insert(user)
	return err
}

// LoadUser get data from a user.
func (mongoDb MongoDBConnection) LoadUser(Email string) (result auth.User, err error) {
	mongoDb.session = mongoDb.GetSession()
	defer mongoDb.session.Close()
	users := mongoDb.session.DB(databaseName).C(usersCollectionKey)
	err = users.Find(bson.M{"email": Email}).One(&result)
	return result, err
}

// GetSession return a new session if there is no previous one.
func (mongoDb *MongoDBConnection) GetSession() *mgo.Session {
	if mongoDb.session != nil {
		return mongoDb.session.Copy()
	}
	session, err := mgo.Dial(databaseHostname)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}
