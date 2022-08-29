package tests

import (
	"github.com/pjmd89/goutils/dbutils"
	"github.com/pjmd89/mongomodel/mongomodel"
)

func OnDB(currentDB string, currentCollection string) (r string) {
	r = currentDB
	if currentCollection == "esta coleccion cambia la BD" {
		r = "otra base de datos"
	}
	r = "tests"
	return r
}

func CreateDB() dbutils.DBInterface {
	var configPath string = "../../etc/db/db.json"

	db := mongomodel.NewConn(&configPath)
	db.(*mongomodel.MongoDBConn).OnDatabase = OnDB

	return db
}
