package mongomodel

import (
	"github.com/pjmd89/goutils/dbutils"
	"go.mongodb.org/mongo-driver/mongo"
)
type MongoDBConnInterface interface {
	GetClient() *mongo.Client
}
type MongoUpdate struct{
	Set interface{} `bson:"$set"`
}
type MongoDBConn struct {
	MongoDBConnInterface
	dbutils.DBInterface
	dbutils.DB
	tryingCounter 			int
	client 					*mongo.Client
	database 				string
	collection				string
	SkipCollection 			[]string		`json:"skipCollection"`
}
