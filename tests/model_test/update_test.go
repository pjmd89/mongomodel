package model_test

import (
	"testing"

	"github.com/pjmd89/mongomodel/tests"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestUpdate(t *testing.T) {
	//var userResult []tests.TestTypes
	db := tests.CreateDB()
	Users := tests.TestTypes{}
	Users.Init(tests.TestTypes{}, db)

	// userData := map[string]interface{}{
	// 	"nombre":   "Pablo",
	// 	"apellido": "Munoz",
	// 	"edad":     30,
	// 	"mapWithVal": map[string]any{
	// 		"nombre": "sebastian",
	// 	},
	// 	"map": map[string]any{
	// 		"nombre": "sebastian",
	// 	},
	// 	"struct": map[string]any{"stringOne": "foo", "stringTwo": "bar"},
	// }

	// _, err := Users.Create(userData, nil)
	// if err != nil {
	// 	t.Fatal("Se genero un error al crear el registro: ", err.Error())
	// }
	id, _ := primitive.ObjectIDFromHex("67698200fbc9721d0e5d0083")
	findById := map[string]interface{}{
		"_id": id,
	}
	updateData := map[string]any{
		"structPtrWithVal": map[string]any{"stringOne": "valueOne", "stringTwo": "valueTwo", "intOne": 3},
		"arr": []string{
			"foo",
			"bar",
		},
		"map": map[string]any{
			"nombre": "jose",
		},
		"mapWithVal": map[string]any{
			"nombre": "santiago",
		},
		"struct": map[string]any{"stringOne": "valueOne", "stringTwo": "valueTwo", "intOne": 4},
	}
	opts := []*options.UpdateOptions{}

	result, err := Users.Update(updateData, findById, opts)
	if err != nil {
		t.Fatal("Se genero un error al actualizar el registro: ", err.Error())
	}
	userResult := result.([]tests.TestTypes)
	t.Logf("Fecha de actualizacion: %d", userResult[0].Updated)
	if userResult[0].Updated == 0 {
		t.Log("El campo Created no fue seteado")
	}

}
