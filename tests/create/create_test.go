package create

import (
	"testing"

	"github.com/pjmd89/mongomodel/tests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreate(t *testing.T){
	db := tests.CreateDB()
	Users := tests.Usuarios{};
	Users.Init(tests.Usuarios{}, db);

	userData := map[string]interface{}{
		"nombre": "Pablo",
		"apellido": "Munoz",
		"edad": 30,
	}

	result,err := Users.Create(userData, nil);

	if err!= nil{
		t.Fatal("Se genero un error al crear el registro: ", err.Error())
	}
	var userResult *tests.Usuarios = result.(*tests.Usuarios);

	id := userResult.Id;
	
	t.Logf("Id registrado: %v", id.Hex());
	if !primitive.IsValidObjectID(id.Hex()){
		t.Fatal("No es un ID valido: ", id.Hex());
	}
	where := map[string]interface{}{
		"_id":id,
	}
	result, err = Users.Read(where, nil);
	if err != nil {
		t.Fatal("Se genero un error al hacer la consulta del registro: ", err.Error());
	}
	if len(result.([]tests.Usuarios)) == 0{
		t.Fatal("No se inserto ningun registro en la base de datos");
	}
}