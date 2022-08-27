package read

import (
	"testing"

	"github.com/pjmd89/mongomodel/tests"
)

func TestRead(t *testing.T){
	var userResult []tests.Usuarios;
	insertValues := 10;
	db := tests.CreateDB()
	Users := tests.Usuarios{};
	Users.Init(tests.Usuarios{}, db);
	
	insertNewData(insertValues);

	result,err := Users.Read(nil,nil);

	if err != nil{
		t.Fatal("Se genero un error al consultar la base de datos:", err.Error())
	}
	userResult = result.([]tests.Usuarios);
	if len(userResult) < insertValues {
		t.Fatalf("el total de registros no es correcto. Total de registros: %v",len(userResult));
	}

	findByID := map[string]interface{}{
		"_id":userResult[0].Id,
	}
	result, err = Users.Read(findByID,nil);
}
func TestReadOne(t *testing.T){
	var userResult []tests.Usuarios;
	db := tests.CreateDB()
	Users := tests.Usuarios{};
	Users.Init(tests.Usuarios{}, db);
	result,err := Users.Read(nil,nil);
	if err != nil{
		t.Fatal("Se genero un error al consultar la base de datos:", err.Error())
	}
	userResult = result.([]tests.Usuarios);
	findByID := map[string]interface{}{
		"_id":userResult[0].Id,
	}
	t.Logf("ID que se esta buscando: %v",userResult[0].Id.Hex());

	result, err = Users.Read(findByID,nil);
	userResult = result.([]tests.Usuarios);
	if err != nil{
		t.Fatal("Se genero un error al consultar la base de datos:", err.Error())
	}

	if len(userResult) > 1{
		t.Fatal("Se encontraron muchos registros");
	}

	if len(userResult) == 0{
		t.Fatal("No se consiguio ningun registro");
	}
}
func insertNewData(values int){
	db := tests.CreateDB()
	Users := tests.Usuarios{};
	Users.Init(tests.Usuarios{}, db);

	userData := map[string]interface{}{
		"nombre": "Pablo",
		"apellido": "Munoz",
		"edad": 30,
	}

	for i:=0; i <values; i++{
		userData["edad"] = userData["edad"].(int)+i; 
		Users.Create(userData, nil);
	}
}