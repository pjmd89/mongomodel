package model_test

import (
	"testing"
	"time"

	"github.com/pjmd89/mongomodel/tests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateID(t *testing.T) {
	db := tests.CreateDB()
	testTypes := tests.TestTypes{}
	testTypes.Init(tests.TestTypes{}, db)

	userData := map[string]interface{}{}

	result, err := testTypes.Create(userData, nil)

	if err != nil {
		t.Fatal("Se genero un error al crear el registro: ", err.Error())
	}
	var testResult *tests.TestTypes = result.(*tests.TestTypes)

	id := testResult.Id
	t.Logf("Id registrado: %v", id.Hex())
	if !primitive.IsValidObjectID(id.Hex()) {
		t.Fatal("No es un ID valido: ", id.Hex())
	}
}
func TestCreateInt(t *testing.T) {
	db := tests.CreateDB()
	saveInt := 20
	Users := tests.TestTypes{}
	Users.Init(tests.TestTypes{}, db)

	userData := map[string]interface{}{
		"intWithVal":    saveInt,
		"intPtrWithVal": saveInt,
	}

	result, err := Users.Create(userData, nil)

	if err != nil {
		t.Fatal("Se genero un error al crear el registro: ", err.Error())
	}
	var testResult *tests.TestTypes = result.(*tests.TestTypes)

	t.Logf("saveInt: %d, IntWithVal: %d", saveInt, testResult.IntWithVal)
	if testResult.IntWithVal != saveInt {
		t.Fatal("el campo IntWithVal no se guardo correctamente")
	}
	t.Logf("saveInt: %d, IntPtrWithVal: %d", saveInt, testResult.IntPtrWithVal)
	if *testResult.IntPtrWithVal != saveInt {
		t.Fatal("el campo IntPtrWithVal no se guardo correctamente")
	}
	t.Logf("Int: %d", testResult.Int)
	if testResult.Int != 0 {
		t.Fatal("el campo Int no se guardo correctamente")
	}
	t.Logf("IntPtr: %d", testResult.IntPtr)
	if testResult.IntPtr != nil {
		t.Fatal("el campo IntPtr no se guardo correctamente")
	}
	t.Logf("IntDef: %d", testResult.IntDef)
	if testResult.IntDef != 120 {
		t.Fatal("el campo IntDef no se guardo correctamente")
	}
	t.Logf("IntPtrDef: %d", testResult.IntPtrDef)
	if *testResult.IntPtrDef != 15 {
		t.Fatal("el campo IntPtrDef no se guardo correctamente")
	}
}
func TestCreateString(t *testing.T) {
	db := tests.CreateDB()
	saveStr := "tests"
	Users := tests.TestTypes{}
	Users.Init(tests.TestTypes{}, db)

	userData := map[string]interface{}{
		"StringWithVal":    saveStr,
		"StringPtrWithVal": saveStr,
	}

	result, err := Users.Create(userData, nil)

	if err != nil {
		t.Fatal("Se genero un error al crear el registro: ", err.Error())
	}
	var testResult *tests.TestTypes = result.(*tests.TestTypes)
	t.Logf("saveStr: %v, StringWithVal: %v", saveStr, testResult.StringWithVal)
	if testResult.StringWithVal != saveStr {
		t.Fatal("el campo StringWithVal no se guardo correctamente")
	}
	t.Logf("saveStr: %s, StringPtrWithVal: %v", saveStr, *testResult.StringPtrWithVal)
	if *testResult.StringPtrWithVal != saveStr {
		t.Fatal("el campo StringPtrWithVal no se guardo correctamente")
	}
	t.Logf("saveStr: %v", testResult.String)
	if testResult.String != "" {
		t.Fatal("el campo String no se guardo correctamente")
	}
	t.Logf("StringPtr: %v", testResult.StringPtr)
	if testResult.StringPtr != nil {
		t.Fatal("el campo StringPtr no se guardo correctamente")
	}
	t.Logf("StringDef: %v", testResult.StringDef)
	if testResult.StringDef != "test default" {
		t.Fatal("el campo StringDef no se guardo correctamente")
	}
	t.Logf("StringPtrDef: %v", testResult.StringPtrDef)
	if *testResult.StringPtrDef != "test ptr default" {
		t.Fatal("el campo StringPtrDef no se guardo correctamente")
	}
}
func TestCreatedDate(t *testing.T) {
	db := tests.CreateDB()
	Users := tests.TestTypes{}
	Users.Init(tests.TestTypes{}, db)

	userData := map[string]interface{}{}

	result, err := Users.Create(userData, nil)

	if err != nil {
		t.Fatal("Se genero un error al crear el registro: ", err.Error())
	}
	var testResult *tests.TestTypes = result.(*tests.TestTypes)
	now := time.Now().Unix()

	t.Logf("Created: %d, now: %d", testResult.Created, now)

	if testResult.Created != now {
		t.Fatal("los tipos de tipo created no estan funcionando correctamente")
	}
	t.Logf("CreatedPtr: %d, now: %d", *testResult.CreatedPtr, now)
	if *testResult.CreatedPtr != now {
		t.Fatal("los tipos de tipo [pointer] created no estan funcionando correctamente")
	}
}
