package model_test

import (
	"testing"

	"github.com/pjmd89/mongomodel/tests"
)

func TestDelete(t *testing.T) {
	db := tests.CreateDB()
	Users := tests.TestTypes{}
	Users.Init(tests.TestTypes{}, db)

	userData := map[string]interface{}{
		"nombre":   "Pablo",
		"apellido": "Munoz",
		"edad":     30,
	}
	result, err := Users.Create(userData, nil)

	if err != nil {
		t.Fatal("Se genero un error al crear el registro: ", err.Error())
	}
	var userInsertResult *tests.TestTypes = result.(*tests.TestTypes)

	findById := map[string]interface{}{
		"_id": userInsertResult.Id,
	}
	t.Logf("Id a elminiar: %v", userInsertResult.Id.Hex())
	result, err = Users.Read(nil, nil)
	counter := len(result.([]tests.TestTypes))

	t.Logf("Total de registros iniciales: %v", counter)
	Users.Delete(findById, nil)
	result, err = Users.Read(nil, nil)
	counterDeleted := len(result.([]tests.TestTypes))
	t.Logf("Total de registros finales: %v", counterDeleted)
	if counterDeleted+1 != counter {
		t.Fatal("Hubo un error en la eliminacion del registro.")
	}

}
