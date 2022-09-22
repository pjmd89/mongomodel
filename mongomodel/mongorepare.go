package mongomodel

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/pjmd89/goutils/dbutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (o *Model) RepareData(self any, data []bson.M, scalarIdType interface{}) (err error) {
	fmt.Println(reflect.TypeOf(self))
	xd := dbutils.CreateStruct(self, scalarIdType, primitive.ObjectID{}, false)
	rType := reflect.TypeOf(xd)
	for _, v := range data {
		instance := reflect.New(rType)
		for i := 0; i < rType.NumField(); i++ {
			typedField := rType.Field(i)
			gqlTags := dbutils.GetTags(typedField)
			bsonTagString := typedField.Tag.Get("bson")
			tags := strings.Split(bsonTagString, ",")
			if tags[0] == "_id" {
				typedField.Tag = reflect.StructTag("`bson:\"-\"`")
			}
			if tags[0] != "_id" && v[tags[0]] != nil && strings.Trim(bsonTagString, " ") != "-" {
				switch typedField.Type.Kind() {
				case reflect.Struct:
					o.repareStruct(instance, typedField.Name, v[tags[0]])
				case reflect.Ptr:
					o.reparePtr(instance, typedField.Name, v[tags[0]])
				case reflect.Array, reflect.Slice:
					o.repareSlice(instance, typedField.Name, v[tags[0]], gqlTags)
				case reflect.Map:
					o.repareMap(instance, typedField.Name, v[tags[0]])
				case reflect.String:
					o.repareString(instance, typedField.Name, v[tags[0]])
				case reflect.Bool:
					o.repareBool(instance, typedField.Name, v[tags[0]])
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					o.repareInt(instance, typedField.Name, v[tags[0]])
				case reflect.Float32, reflect.Float64:
					o.repareFloat(instance, typedField.Name, v[tags[0]])
				}
			}
		}
		x := instance.Interface()
		///*

		where := map[string]interface{}{
			"_id": v["_id"],
		}
		_, err := o.InterfaceUpdate(x, where, nil)
		if err != nil {
			log.Println(err.Error())
		}
		//*/
	}
	return
}
func (o *Model) repareStruct(value reflect.Value, fieldName string, data any) (r reflect.Value) {

	return
}
func (o *Model) repareString(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	rValue := value.Elem().FieldByName(fieldName)
	var sData reflect.Value
	switch data.(type) {
	case primitive.ObjectID:
		sData = reflect.ValueOf(data)
	default:
		sData = reflect.ValueOf(fmt.Sprintf("%v", data))
	}

	rValue.Set(sData.Convert(rValue.Type()))
	return
}
func (o *Model) repareSlice(value reflect.Value, fieldName string, data any, tags dbutils.Tags) (r reflect.Value) {
	parse := value.Elem().FieldByName(fieldName)
	var sData reflect.Value = reflect.ValueOf(data)
	switch parse.Type() {
	case reflect.TypeOf(primitive.ObjectID{}):
		switch vData := data.(type) {
		case primitive.ObjectID:
			parse.Set(sData)
		case string:
			nId, _ := primitive.ObjectIDFromHex(vData)
			parse.Set(reflect.ValueOf(nId))
		}
	default:
		vData := reflect.ValueOf(data)
		switch parse.Type().Elem() {
		case reflect.TypeOf(primitive.ObjectID{}):
			var idContainers []primitive.ObjectID
			for i := 0; i < vData.Len(); i++ {
				var idData primitive.ObjectID
				switch iData := vData.Index(i).Interface().(type) {
				case primitive.ObjectID:
					idData = iData
				case string:
					idData, _ = primitive.ObjectIDFromHex(iData)
				}
				idContainers = append(idContainers, idData)
			}
		default:
			parse.Set(vData)
		}
	}
	return
}
func (o *Model) reparePtr(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	if !value.IsNil() {

	}
	return
}
func (o *Model) repareMap(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	return
}
func (o *Model) repareBool(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	parse, _ := strconv.ParseBool(fmt.Sprintf("%v", data))
	rValue := value.Elem().FieldByName(fieldName)
	sData := reflect.ValueOf(parse)
	rValue.Set(sData.Convert(rValue.Type()))
	return
}
func (o *Model) repareInt(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	parse, _ := strconv.ParseInt(fmt.Sprintf("%v", data), 10, 64)
	rValue := value.Elem().FieldByName(fieldName)
	sData := reflect.ValueOf(parse)
	rValue.Set(sData.Convert(rValue.Type()))
	return
}
func (o *Model) repareFloat(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	parse, _ := strconv.ParseFloat(fmt.Sprintf("%v", data), 64)
	rValue := value.Elem().FieldByName(fieldName)
	sData := reflect.ValueOf(parse)
	rValue.Set(sData.Convert(rValue.Type()))
	return
}
