package mongomodel

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/pjmd89/goutils/dbutils"
	"github.com/pjmd89/goutils/systemutils/debugmode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (o *Model) RepareData(self any, data []bson.M) (err error) {
	xd := dbutils.CreateStruct(self, false)
	rType := reflect.TypeOf(xd)
	for _, v := range data {
		x := o.repareStruct(rType, v)
		where := map[string]interface{}{
			"_id": v["_id"],
		}
		if !debugmode.Enabled {
			_, err := o.InterfaceReplace(x.Interface(), where, nil)
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			//log.Println(x.Interface())
		}
	}
	return
}
func (o *Model) parseField(typedField reflect.StructField, instance reflect.Value, v bson.M) {
	var realTag string
	var valueData interface{}
	var composeData []string
	gqlTags := dbutils.GetTags(typedField)
	bsonTagString := typedField.Tag.Get("bson")
	tags := strings.Split(bsonTagString, ",")
	realTag = tags[0]
	valueData = v[tags[0]]
	if tags[0] == "_id" {
		typedField.Tag = reflect.StructTag("`bson:\"-\"`")
	}

	if len(gqlTags.Compose) > 0 {
		for _, vv := range gqlTags.Compose {
			vData := fmt.Sprintf("%v", v[vv])
			composeData = append(composeData, vData)
		}
		v[tags[0]] = strings.Trim(strings.Join(composeData, " "), " ")
	}
	if gqlTags.Change != "" {
		realTag = gqlTags.Change
	}
	if v[realTag] != nil && valueData == nil {
		valueData = v[realTag]
	}

	if tags[0] != "_id" && strings.Trim(bsonTagString, " ") != "-" {
		switch typedField.Type.Kind() {
		case reflect.Struct:
			var sendInstance reflect.Type
			sendInstance = instance.Type()
			if instance.Kind() == reflect.Ptr {
				sendInstance = instance.Type().Elem()
			}
			if valueData != nil {
				switch vData := valueData.(type) {
				case bson.M:
					o.repareStruct(sendInstance, vData)
				case bson.A:
					for _, iv := range vData {
						fmt.Println(reflect.ValueOf(iv).Type())
					}
				}
			}

		case reflect.Ptr:
			o.reparePtr(typedField, instance, typedField.Name, valueData)
		case reflect.Array, reflect.Slice:
			o.repareSlice(instance, typedField.Name, valueData, gqlTags, typedField)
		case reflect.Map:
			o.repareMap(instance, typedField.Name, valueData)
		case reflect.String:
			o.repareString(instance, typedField.Name, valueData)
		case reflect.Bool:
			o.repareBool(instance, typedField.Name, valueData)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			o.repareInt(instance, typedField.Name, valueData)
		case reflect.Float32, reflect.Float64:
			o.repareFloat(instance, typedField.Name, valueData)
		default:
			x := typedField.Type.Kind()
			fmt.Println(x)
		}
	}
}
func (o *Model) repareStruct(rType reflect.Type, v bson.M) (r reflect.Value) {
	instance := reflect.New(rType)
	for i := 0; i < rType.NumField(); i++ {
		o.parseField(rType.Field(i), instance, v)
	}
	return instance.Elem()
}
func (o *Model) parsePtr(typedField reflect.StructField, instance reflect.Value, v any) {
	switch instance.Type().Elem().Kind() {
	case reflect.Struct:
		switch v.(type) {
		case bson.M:
			o.repareStruct(instance.Type(), v.(bson.M))
		}
	case reflect.Ptr:

	case reflect.Array, reflect.Slice:
		gqlTags := dbutils.GetTags(typedField)
		o.repareSlice(instance, typedField.Name, v, gqlTags, typedField)
	case reflect.Map:
		o.repareMap(instance, typedField.Name, v)
	case reflect.String:
		o.repareString(instance, typedField.Name, v)
	case reflect.Bool:
		o.repareBool(instance, typedField.Name, v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		o.repareInt(instance, typedField.Name, v)
	case reflect.Float32, reflect.Float64:
		o.repareFloat(instance, typedField.Name, v)
	}
}
func (o *Model) repareString(value reflect.Value, fieldName string, data any) (r reflect.Value) {

	if value.IsNil() {
		return
	}
	if data == nil {
		switch value.Elem().FieldByName(fieldName).Type().Kind() {
		case reflect.String:
			x := ""
			value.Elem().FieldByName(fieldName).Set(reflect.ValueOf(x).Convert(value.Elem().FieldByName(fieldName).Type()))
		}
		return
	}
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
func (o *Model) repareSlice(value reflect.Value, fieldName string, data any, tags dbutils.Tags, typedField reflect.StructField) (r reflect.Value) {
	parse := value
	if data == nil {
		switch typedField.Type.Kind() {
		case reflect.Ptr:
		case reflect.Array, reflect.Slice:
			if typedField.Type == reflect.TypeOf(primitive.ObjectID{}) {
				x := reflect.ValueOf(primitive.ObjectID{})
				parse.Elem().FieldByName(fieldName).Set(x)
			} else {
				x := reflect.MakeSlice(reflect.SliceOf(typedField.Type.Elem()), 0, 0)
				parse.Elem().FieldByName(fieldName).Set(x)
			}

		}
		return
	}
	isValid := value.Elem().IsValid()
	if isValid {
		parse = value.Elem().FieldByName(fieldName)
	}

	var sData reflect.Value = reflect.ValueOf(data)

	switch parse.Interface().(type) {
	case primitive.ObjectID:
		switch vxData := data.(type) {
		case primitive.ObjectID:
			parse.Set(sData)
		case string:
			nId, _ := primitive.ObjectIDFromHex(vxData)
			parse.Set(reflect.ValueOf(nId))
		}
	case []primitive.ObjectID:
		x := reflect.MakeSlice(reflect.SliceOf(typedField.Type.Elem()), 0, 0)
		for _, id := range data.(primitive.A) {
			switch id.(type) {
			case string:
				newObjectID, _ := primitive.ObjectIDFromHex(id.(string))
				x = reflect.Append(x, reflect.ValueOf(newObjectID))
			case primitive.ObjectID:
				x = reflect.Append(x, reflect.ValueOf(id))
			}

		}
		parse.Set(x)
	case *primitive.ObjectID:
		x := reflect.New(reflect.TypeOf(primitive.ObjectID{}))
		switch vxData := data.(type) {
		case primitive.ObjectID:
			x.Elem().Set(sData)
			parse.Set(x)
		case string:
			nId, _ := primitive.ObjectIDFromHex(vxData)
			x.Elem().Set(reflect.ValueOf(nId))
			parse.Set(x)
		}
	default:
		//vData := reflect.ValueOf(data)
		newx := reflect.New(parse.Type().Elem()).Elem()
		switch newx.Interface().(type) {
		case primitive.ObjectID:
			var idContainers []primitive.ObjectID
			if data != nil {
				for i := 0; i < sData.Len(); i++ {
					var idData primitive.ObjectID
					switch iData := sData.Index(i).Interface().(type) {
					case primitive.ObjectID:
						idData = iData
					case string:
						idData, _ = primitive.ObjectIDFromHex(iData)
					}
					idContainers = append(idContainers, idData)
				}
			}
			parse.Set(reflect.ValueOf(idContainers))
		case string:
			iContainers := make([]string, 0, 0)
			if fieldName == "MemoryTags" {
				fmt.Println(sData)
			}
			if data != nil && !sData.IsNil() {
				for i := 0; i < sData.Len(); i++ {
					iContainers = append(iContainers, sData.Index(i).Interface().(string))
				}
			}
			parse.Set(reflect.ValueOf(iContainers))
		default:
			count := 0
			if data != nil && !sData.IsNil() {
				count = sData.Len()
			}
			if count == 0 {
				x := reflect.MakeSlice(reflect.SliceOf(parse.Type().Elem()), 0, 0)
				parse.Set(x)
			}
			for i := 0; i < count; i++ {
				x := parse.Type().Elem().Kind()
				switch x {
				case reflect.Struct:
					parse.Set(reflect.Append(parse, o.repareStruct(parse.Type().Elem(), sData.Index(i).Interface().(bson.M))))
				default:
					parse.Set(reflect.Append(parse, sData.Index(i)))
				}
			}
		}
	}
	return
}
func (o *Model) reparePtr(typedField reflect.StructField, value reflect.Value, fieldName string, data any) (r reflect.Value) {
	if data != nil {
		o.parsePtr(typedField, value.Elem().FieldByName(fieldName), data)
	}
	return
}
func (o *Model) repareMap(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	return
}
func (o *Model) repareBool(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	if data == nil {
		switch value.Elem().FieldByName(fieldName).Type().Kind() {
		case reflect.Bool:
			x := false
			value.Elem().FieldByName(fieldName).Set(reflect.ValueOf(x).Convert(value.Elem().FieldByName(fieldName).Type()))
		}
		return
	}
	parse, _ := strconv.ParseBool(fmt.Sprintf("%v", data))
	rValue := value.Elem().FieldByName(fieldName)
	sData := reflect.ValueOf(parse)
	rValue.Set(sData.Convert(rValue.Type()))
	return
}
func (o *Model) repareInt(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	if data == nil {
		switch value.Elem().FieldByName(fieldName).Type().Kind() {
		case reflect.Int64:
			x := int64(0)
			value.Elem().FieldByName(fieldName).Set(reflect.ValueOf(x).Convert(value.Elem().FieldByName(fieldName).Type()))
		}
		return
	}
	parse, _ := strconv.ParseInt(fmt.Sprintf("%v", data), 10, 64)
	rValue := value.Elem().FieldByName(fieldName)
	sData := reflect.ValueOf(parse)
	rValue.Set(sData.Convert(rValue.Type()))
	return
}
func (o *Model) repareFloat(value reflect.Value, fieldName string, data any) (r reflect.Value) {
	if data == nil {
		switch value.Elem().FieldByName(fieldName).Type().Kind() {
		case reflect.Float64:
			x := float64(0)
			value.Elem().FieldByName(fieldName).Set(reflect.ValueOf(x).Convert(value.Elem().FieldByName(fieldName).Type()))
		}
		return
	}
	parse, _ := strconv.ParseFloat(fmt.Sprintf("%v", data), 64)
	rValue := value.Elem().FieldByName(fieldName)
	sData := reflect.ValueOf(parse)
	rValue.Set(sData.Convert(rValue.Type()))
	return
}
