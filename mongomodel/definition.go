package mongomodel

import (
	"github.com/pjmd89/goutils/dbutils"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBConnInterface interface {
	GetClient() *mongo.Client
}
type Update struct {
	Set interface{} `bson:"$set"`
}
type DatesController struct {
	Created *int64
	Updated *int64
}
type MongoDBConn struct {
	MongoDBConnInterface
	dbutils.DBInterface
	dbutils.DB
	tryingCounter int
	client        *mongo.Client
	database      string
	collection    string
}

type URIData struct {
	Host     string
	Port     string
	User     string
	Pass     string
	DataBase string
}
