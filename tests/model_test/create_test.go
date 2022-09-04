package model_test

import (
	"testing"
	"time"

	"github.com/pjmd89/mongomodel/tests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestID(t *testing.T) {
	db := tests.CreateDB()
	testTypes := tests.TestTypes{}
	testTypes.Init(tests.TestTypes{}, db)
	strID := "62c68a2aed60c36f6253251d"
	objectID, _ := primitive.ObjectIDFromHex(strID)
	userData := map[string]interface{}{
		"idWithVal":      strID,
		"idPtrWithVal":   strID,
		"idWithIDVal":    objectID,
		"idPtrWithIDVal": objectID,
	}

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
	t.Logf("strID: %s, IdWithVal: %s", strID, testResult.IdWithVal.Hex())
	if strID != testResult.IdWithVal.Hex() {
		t.Fatal("IdWithVal no se guardo correctamente")
	}
	t.Logf("strID: %s, IdPtrWithVal: %s", strID, testResult.IdPtrWithVal.Hex())
	if strID != testResult.IdPtrWithVal.Hex() {
		t.Fatal("IdPtrWithVal no se guardo correctamente")
	}
	t.Logf("strID: %s, IdWithIDVal: %s", strID, testResult.IdWithIDVal.Hex())
	if strID != testResult.IdWithIDVal.Hex() {
		t.Fatal("IdWithVal no se guardo correctamente")
	}
	t.Logf("strID: %s, IdPtrWithIDVal: %s", strID, testResult.IdPtrWithIDVal.Hex())
	if strID != testResult.IdPtrWithIDVal.Hex() {
		t.Fatal("IdPtrWithIDVal no se guardo correctamente")
	}
	t.Logf("IdOutVal: %s", testResult.IdOutVal.Hex())
	if "000000000000000000000000" != testResult.IdOutVal.Hex() {
		t.Fatal("IdOutVal no se guardo correctamente")
	}
	t.Logf("IdPtrOutVal: %v", testResult.IdPtrOutVal)
	if testResult.IdPtrOutVal != nil {
		t.Fatal("IdPtrOutVal no se guardo correctamente")
	}
}
func TestInt(t *testing.T) {
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
func TestString(t *testing.T) {
	db := tests.CreateDB()
	saveStr := "tests"
	Users := tests.TestTypes{}
	Users.Init(tests.TestTypes{}, db)

	userData := map[string]interface{}{
		"stringWithVal":    saveStr,
		"stringPtrWithVal": saveStr,
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
func TestArray(t *testing.T) {
	db := tests.CreateDB()
	Users := tests.TestTypes{}
	Users.Init(tests.TestTypes{}, db)
	arrData := []string{
		"Hola",
		"Mundo",
	}

	arrStructData := []interface{}{
		map[string]interface{}{
			"stringOne": "string",
		},
	}
	userData := map[string]interface{}{
		"arrWithVal":          arrData,
		"arrPtrWithVal":       arrData,
		"arrStructWithVal":    arrStructData,
		"arrPtrStructWithVal": arrStructData,
	}

	result, err := Users.Create(userData, nil)

	if err != nil {
		t.Fatal("Se genero un error al crear el registro: ", err.Error())
	}
	var testResult *tests.TestTypes = result.(*tests.TestTypes)

	t.Logf("Arr: %v", testResult.Arr)
	if testResult.Arr == nil {
		t.Fatal("El campo Arr no se guardo correctamente")
	}
	t.Logf("ArrPtr: %v", testResult.ArrPtr)
	if testResult.ArrPtr != nil {
		t.Fatal("El campo ArrPtr no se guardo correctamente")
	}
	t.Logf("ArrWithVal: %v", testResult.ArrWithVal)
	if testResult.ArrWithVal[0] != "Hola" {
		t.Fatal("El campo Arr no se guardo correctamente")
	}
	t.Logf("ArrPtrWithVal: %v", testResult.ArrPtrWithVal)
	if (*testResult.ArrPtrWithVal)[0] != "Hola" {
		t.Fatal("El campo ArrPtr no se guardo correctamente")
	}
	t.Logf("ArrStruct: %v", testResult.ArrStruct)
	if len(testResult.ArrStruct) != 0 {
		t.Fatal("El campo ArrStruct no se guardo correctamente")
	}
	t.Logf("ArrPtrStruct: %v", testResult.ArrPtrStruct)
	if testResult.ArrPtrStruct != nil {
		t.Fatal("El campo ArrPtrStruct no se guardo correctamente")
	}

	t.Logf("ArrStructWithVal: %v", testResult.ArrStructWithVal)
	if len(testResult.ArrStructWithVal) == 0 {
		t.Fatal("El campo ArrStructWithVal no se guardo correctamente")
	}
	if testResult.ArrStructWithVal[0].StringOne != "string" {
		t.Fatal("El campo ArrStructWithVal no se guardo correctamente")
	}
	t.Logf("ArrPtrStructWithVal: %v", testResult.ArrPtrStructWithVal)
	if testResult.ArrPtrStructWithVal == nil {
		t.Fatal("El campo ArrPtrStructWithVal no se guardo correctamente")
	}
	if (*testResult.ArrPtrStructWithVal)[0].StringOne != "string" {
		t.Fatal("El campo ArrPtrStructWithVal no se guardo correctamente")
	}
}
func TestStruct(t *testing.T) {
	db := tests.CreateDB()
	Users := tests.TestTypes{}
	Users.Init(tests.TestTypes{}, db)
	structata := map[string]interface{}{
		"stringOne": "str",
	}
	userData := map[string]interface{}{
		"structWithVal":    structata,
		"structPtrWithVal": structata,
	}

	result, err := Users.Create(userData, nil)
	if err != nil {
		t.Fatal("Se genero un error al crear el registro: ", err.Error())
	}
	var testResult *tests.TestTypes = result.(*tests.TestTypes)

	t.Logf("Struct: %v", testResult.Struct)
	if testResult.Struct.StringOne != "" {
		t.Fatal("El campo Struct no se guardo correctamente")
	}
	if testResult.Struct.StringTwo != "string2" {
		t.Fatal("El campo Struct no se guardo correctamente")
	}
	t.Logf("StructPtr: %v", testResult.StructPtr)
	if testResult.StructPtr != nil {
		t.Fatal("El campo StructPtr no se guardo correctamente")
	}

	t.Logf("StructWithVal: %v", testResult.StructWithVal)
	if testResult.StructWithVal.StringOne != "str" {
		t.Fatal("El campo StructWithVal no se guardo correctamente")
	}
	if testResult.StructWithVal.StringTwo != "string2" {
		t.Fatal("El campo StructWithVal no se guardo correctamente")
	}
	t.Logf("StructPtrWithVal: %v", testResult.StructPtrWithVal)
	if testResult.StructPtrWithVal == nil {
		t.Fatal("El campo StructPtrWithVal no se guardo correctamente")
	}
	if testResult.StructPtrWithVal.StringOne != "str" {
		t.Fatal("El campo StructPtrWithVal no se guardo correctamente")
	}
	if testResult.StructPtrWithVal.StringTwo != "string2" {
		t.Fatal("El campo StructPtrWithVal no se guardo correctamente")
	}
}
func TestDate(t *testing.T) {
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
