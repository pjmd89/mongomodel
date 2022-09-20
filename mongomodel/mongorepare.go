package mongomodel

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (o *Model) RepareData(self any, data []bson.M) (err error) {
	rType := reflect.TypeOf(self)

	for _, v := range data {
		instance := reflect.New(rType)
		for i := 0; i < rType.NumField(); i++ {
			typedField := rType.Field(i)
			//tags := dbutils.GetTags(typedField)
			bsonTagString := typedField.Tag.Get("bson")
			tags := strings.Split(bsonTagString, ",")
			if v[tags[0]] != nil && strings.Trim(bsonTagString, " ") != "-" {
				switch typedField.Type.Kind() {
				case reflect.Struct:
					o.repareStruct(instance, typedField.Name, v[tags[0]])
				case reflect.Ptr:
					o.reparePtr(instance, typedField.Name, v[tags[0]])
				case reflect.Array, reflect.Slice:
					o.repareSlice(instance, typedField.Name, v[tags[0]])
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
		where := map[string]interface{}{
			"_id": v["_id"],
		}
		_, err := o.InterfaceUpdate(x, where, nil)
		if err != nil {
			log.Println(err.Error())
		}
		//fmt.Println(x)
	}
	return
}
func (o *Model) repareStruct(value reflect.Value, fieldName string, data any) (r reflect.Value) {

	return
}
func (o *Model) repareString(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	rValue := value.Elem().FieldByName(fieldName)
	sData := reflect.ValueOf(fmt.Sprintf("%v", data))
	rValue.Set(sData.Convert(rValue.Type()))
	return
}
func (o *Model) repareSlice(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	parse := value.Elem().FieldByName(fieldName)

	switch parse.Type() {
	case reflect.TypeOf(primitive.ObjectID{}):
		dataVal := fmt.Sprintf("%v", data)
		id, _ := primitive.ObjectIDFromHex(dataVal)
		parse.Set(reflect.ValueOf(id))
		return
	default:
		fmt.Println("")
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
