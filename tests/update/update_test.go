package update

import (
	"testing"

	"github.com/pjmd89/mongomodel/tests"
)

func TestUpdate(t *testing.T) {
	var userResult []tests.Usuarios
	db := tests.CreateDB()
	Users := tests.Usuarios{}
	Users.Init(tests.Usuarios{}, db)

	userData := map[string]interface{}{
		"nombre":   "Pablo",
		"apellido": "Munoz",
		"edad":     30,
	}

	result, err := Users.Create(userData, nil)

	if err != nil {
		t.Fatal("Se genero un error al crear el registro: ", err.Error())
	}
	var userInsertResult *tests.Usuarios = result.(*tests.Usuarios)

	if userInsertResult.Edad != 30 {
		t.Fatal("Hay un error en la edad, verifica si se guardo correctamente el registro.")
	}
	findById := map[string]interface{}{
		"_id": userInsertResult.Id,
	}
	updateData := map[string]interface{}{
		"edad": 34,
	}
	
	result, err = Users.Update(updateData, findById, nil)
	userResult = result.([]tests.Usuarios)

	t.Logf("Id a actualizar: %v", userInsertResult.Id.Hex())

	if len(userResult) != 1 {
		t.Fatal("No se actualizo correctamente el registro.")
	}
	t.Logf("Fecha de actualizacion: %d", userResult[0].Updated)
	if userResult[0].Updated == nil {
		t.Log("El campo Created no fue seteado")
	}
	if userResult[0].Edad != 34 {
		t.Fatal("Hubo un error al actualizar el registro.")
	}

}
