package model_test

import (
	"fmt"
	"testing"

	"github.com/pjmd89/mongomodel/tests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRead(t *testing.T) {
	var userResult []tests.TestTypes
	insertValues := 10
	db := tests.CreateDB()
	Users := tests.TestTypes{}
	Users.Init(tests.TestTypes{}, db)

	insertNewData(insertValues)

	result, err := Users.Read(nil, nil)

	if err != nil {
		t.Fatal("Se genero un error al consultar la base de datos:", err.Error())
	}
	userResult = result.([]tests.TestTypes)
	if len(userResult) < insertValues {
		t.Fatalf("el total de registros no es correcto. Total de registros: %v", len(userResult))
	}

	findByID := map[string]interface{}{
		"_id": userResult[0].Id,
	}
	result, err = Users.Read(findByID, nil)
}
func TestReadOne(t *testing.T) {
	var userResult []tests.TestTypes
	db := tests.CreateDB()
	Users := tests.TestTypes{}
	Users.Init(tests.TestTypes{}, db)
	id, _ := primitive.ObjectIDFromHex("67698200fbc9721d0e5d0083")
	findByID := map[string]interface{}{
		"_id": id,
	}

	result, err := Users.Read(findByID, nil)
	if err != nil {
		t.Fatal("Se genero un error al consultar la base de datos:", err.Error())
	}
	userResult = result.([]tests.TestTypes)
	if len(userResult) == 0 {
		t.Fatal("No se consiguio ningun registro")
	}
	fmt.Printf("%+v", userResult[0])

}
func insertNewData(values int) {
	db := tests.CreateDB()
	Users := tests.TestTypes{}
	Users.Init(tests.TestTypes{}, db)

	userData := map[string]interface{}{
		"nombre":   "Pablo",
		"apellido": "Munoz",
		"edad":     30,
	}

	for i := 0; i < values; i++ {
		userData["edad"] = userData["edad"].(int) + i
		Users.Create(userData, nil)
	}
}
