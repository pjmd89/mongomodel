package model

import (
	"fmt"
	"testing"

	"github.com/pjmd89/mongomodel/tests"
)

func wellPrint(m []map[string]any) {
	for i := 0; i < len(m); i++ {
		fmt.Println("--------------------------document ", i+1, "--------------------------")
		for k, v := range m[i] {
			fmt.Printf("%s: %s\n", k, v)
		}
	}
}
func TestRepareData(t *testing.T) {
	db := tests.CreateDB()
	testTypes := tests.TestTypes{}
	testTypes.Init(tests.TestTypes{}, db)

	_, err := testTypes.RawRead(nil, nil)
	if err != nil {
		t.Fatal("Se genero un error al leer los registros la colección: ", err.Error())
	}

	//wellPrint(r.([]map[string]any))

	_, err = testTypes.Repare()
	if err != nil {
		t.Fatal("Se genero un error al reparar la colección: ", err.Error())
	}

	_, err = testTypes.RawRead(nil, nil)
	if err != nil {
		t.Fatal("Se genero un error al leer los registros la colección: ", err.Error())
	}

	//wellPrint(r.([]map[string]any))
}
