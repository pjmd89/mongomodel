package delete

import (
	"testing"

	"github.com/pjmd89/mongomodel/tests"
)

func TestDelete(t *testing.T){
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
	var userInsertResult *tests.Usuarios = result.(*tests.Usuarios);
	
	if userInsertResult.Edad != 30{
		t.Fatal("Hay un error en la edad, verifica si se guardo correctamente el registro.");
	}
	
	findById := map[string]interface{}{
		"_id":userInsertResult.Id,
	}
	t.Logf("Id a elminiar: %v",userInsertResult.Id.Hex())
	result,err = Users.Read(nil,nil);
	counter := len(result.([]tests.Usuarios));
	
	t.Logf("Total de registros iniciales: %v",counter)
	Users.Delete(findById,nil);
	result,err = Users.Read(nil,nil);
	counterDeleted := len(result.([]tests.Usuarios));
	t.Logf("Total de registros finales: %v", counterDeleted)
	if counterDeleted + 1 != counter {
		t.Fatal("Hubo un error en la eliminacion del registro.")
	}
	
}