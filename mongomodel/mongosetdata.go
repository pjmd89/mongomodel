package mongomodel

import (
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetID(model interface{}, id primitive.ObjectID) (err error) {
	modelElem := reflect.ValueOf(model)
	isPtr := false
	err = fmt.Errorf("model is not a pointer")
	if modelElem.Kind() == reflect.Ptr {
		modelElem = modelElem.Elem()
		isPtr = true
		err = nil
	}

	if isPtr {
		for i := 0; i < modelElem.NumField(); i++ {
			field := modelElem.Field(i)
			tags := getTags(modelElem.Type().Field(i))
			if tags.IsID {
				field.Set(reflect.ValueOf(id))
				break
			}
		}
	}
	return err
}
func SetData(inputs map[string]interface{}, model interface{}) (r interface{}) {
	modelType := reflect.TypeOf(model)
	modelKind := modelType.Kind()

	switch modelKind {
	case reflect.Struct:
		r = setStruct(inputs, model)
	case reflect.Ptr:
		r = SetData(inputs, reflect.ValueOf(model).Elem().Interface())
	}
	return r
}
func SetCreatedDate(instance interface{}) (err error) {
	valueOf := reflect.ValueOf(instance)
	for i := 0; i < valueOf.Elem().NumField(); i++ {
		field := valueOf.Elem().Field(i)
		fieldType := valueOf.Elem().Type().Field(i)
		fieldKind := field.Kind()
		tag := getTags(fieldType)
		if tag.CreatedDate {
			switch fieldKind {
			case reflect.Ptr:
				fKind := field.Type().Elem().Kind()
				switch fKind {
				case reflect.Int64:
					date := time.Now().Unix()
					newDate := reflect.New(reflect.TypeOf(date))
					newDate.Elem().SetInt(date)
					field.Set(newDate)
					break
				default:
					err = fmt.Errorf("The field %v is not an *int64 type", field.Type().Name())
				}
			default:
				err = fmt.Errorf("The field %v is not a Pointer", field.Type().Name())
			}
		}
	}
	return err
}
func SetUpdatedDate(instance interface{}) (err error) {
	valueOf := reflect.ValueOf(instance)
	for i := 0; i < valueOf.Elem().NumField(); i++ {
		field := valueOf.Elem().Field(i)
		fieldType := valueOf.Elem().Type().Field(i)
		fieldKind := field.Kind()
		tag := getTags(fieldType)
		if tag.UpdatedDate {
			switch fieldKind {
			case reflect.Ptr:
				fKind := field.Type().Elem().Kind()
				switch fKind {
				case reflect.Int64:
					date := time.Now().Unix()
					newDate := reflect.New(reflect.TypeOf(date))
					newDate.Elem().SetInt(date)
					field.Set(newDate)
					break
				default:
					err = fmt.Errorf("The field %v is not an *int64 type", field.Type().Name())
				}
			default:
				err = fmt.Errorf("The field %v is not a Pointer", field.Type().Name())
			}
		}
	}
	return err
}
func setStruct(inputs map[string]interface{}, model interface{}) interface{} {
	newModel := reflect.New(reflect.TypeOf(model))
	for i := 0; i < newModel.Elem().NumField(); i++ {
		field := newModel.Elem().Field(i)
		fieldType := newModel.Elem().Type().Field(i)
		fieldKind := field.Type().Kind()
		tag := getTags(fieldType)
		if inputs[tag.Name] != nil {
			switch fieldKind {
			case reflect.Struct, reflect.Ptr:
				inputType := reflect.TypeOf(inputs[tag.Name])
				inputKind := inputType.Kind()
				switch inputKind {
				case reflect.Struct, reflect.Ptr:
					rField := setStruct(inputs[tag.Name].(map[string]interface{}), field.Interface())
					field.Set(reflect.ValueOf(rField))
				}
				break
			case reflect.Slice, reflect.Array:
				if tag.IsObjectID {
					newID, _ := primitive.ObjectIDFromHex(inputs[tag.Name].(string))
					field.Set(reflect.ValueOf(newID))
				}
				break
			case reflect.Map:
				break
			case reflect.String:
				field.Set(reflect.ValueOf(inputs[tag.Name].(string)))
				break
			case reflect.Int:
				field.Set(reflect.ValueOf(inputs[tag.Name].(int)))
				break
			case reflect.Int8:
				field.Set(reflect.ValueOf(inputs[tag.Name].(int8)))
				break
			case reflect.Int16:
				field.Set(reflect.ValueOf(inputs[tag.Name].(int16)))
				break
			case reflect.Int32:
				field.Set(reflect.ValueOf(inputs[tag.Name].(int32)))
				break
			case reflect.Int64:
				field.Set(reflect.ValueOf(inputs[tag.Name].(int64)))
				break
			case reflect.Float32:
				field.Set(reflect.ValueOf(inputs[tag.Name].(float32)))
				break
			case reflect.Float64:
				field.Set(reflect.ValueOf(inputs[tag.Name].(float64)))
				break
			case reflect.Bool:
				field.Set(reflect.ValueOf(inputs[tag.Name].(bool)))
				break
			}
		}
	}
	return newModel.Interface()
}
func setField() {

}
