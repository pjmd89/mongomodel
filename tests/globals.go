package tests

import (
	"github.com/pjmd89/goutils/dbutils"
	"github.com/pjmd89/mongomodel/mongomodel"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trabajo struct {
	Nombre    string `bson:"nombre"`
	Direccion string `bson:"direccion"`
}
type Usuarios struct {
	mongomodel.Model `bson:"-"`
	Id               primitive.ObjectID `bson:"_id,omitempty" gql:"name=id,id=true,objectID=true"`
	Nombre           string             `bson:"nombre" gql:"name=nombre"`
	Apellido         string             `bson:"apellido" gql:"name=apellido"`
	Edad             int                `bson:"edad" gql:"name=edad"`
	Trabajos         []Trabajo          `bson:"trabajos" gql:"name=trabajos"`
	Created          *int64             `bson:"created" gql:"name=created,createdDate=true"`
	Updated          *int64             `bson:"updated" gql:"name=updated,updatedDate=true"`
}

func OnDB(currentDB string, currentCollection string) (r string) {
	r = currentDB
	if currentCollection == "esta coleccion cambia la BD" {
		r = "otra base de datos"
	}
	if currentCollection == "Usuarios" {
		r = "Pruebita"
	}
	return r
}

func CreateDB() dbutils.DBInterface {
	var configPath string = "../../etc/db/db.json"

	db := mongomodel.NewConn(&configPath)
	db.(*mongomodel.MongoDBConn).OnDatabase = OnDB

	return db
}
