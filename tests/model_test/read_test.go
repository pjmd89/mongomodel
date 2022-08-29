package model_test

import (
	"testing"

	"github.com/pjmd89/mongomodel/tests"
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
	result, err := Users.Read(nil, nil)
	if err != nil {
		t.Fatal("Se genero un error al consultar la base de datos:", err.Error())
	}
	userResult = result.([]tests.TestTypes)
	findByID := map[string]interface{}{
		"_id": userResult[0].Id,
	}
	t.Logf("ID que se esta buscando: %v", userResult[0].Id.Hex())

	result, err = Users.Read(findByID, nil)
	userResult = result.([]tests.TestTypes)
	if err != nil {
		t.Fatal("Se genero un error al consultar la base de datos:", err.Error())
	}

	if len(userResult) > 1 {
		t.Fatal("Se encontraron muchos registros")
	}

	if len(userResult) == 0 {
		t.Fatal("No se consiguio ningun registro")
	}
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
