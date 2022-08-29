package mongomodel

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

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
func SetData(inputs map[string]interface{}, model interface{}, datesController DatesController) (r interface{}) {
	modelType := reflect.TypeOf(model)
	modelKind := modelType.Kind()

	switch modelKind {
	case reflect.Struct:
		r = setStruct(inputs, model, datesController)
	case reflect.Ptr:
		r = SetData(inputs, reflect.ValueOf(model).Elem().Interface(), datesController)
	}
	return r
}
func setStruct(inputs map[string]interface{}, model interface{}, datesController DatesController) interface{} {
	newModel := reflect.New(reflect.TypeOf(model))
	for i := 0; i < newModel.Elem().NumField(); i++ {
		field := newModel.Elem().Field(i)
		fieldType := newModel.Elem().Type().Field(i)
		fieldKind := field.Type().Kind()
		tag := getTags(fieldType)
		if tag.CreatedDate && datesController.Created != nil {
			inputs[tag.Name] = *datesController.Created
		}
		if tag.UpdatedDate && datesController.Updated != nil {
			inputs[tag.Name] = *datesController.Updated
		}
		if strings.Trim(tag.Name, " ") != "" {
			if inputs[tag.Name] != nil {
				setDataOn(inputs, tag, fieldKind, field, datesController)
			} else {
				setNilOn(tag, fieldKind, field)
			}
		}
	}
	return newModel.Interface()
}
func setNilOn(tag Tags, fieldKind reflect.Kind, field reflect.Value) {
	switch fieldKind {
	case reflect.Struct:
		//fmt.Println(tag.Name, fieldKind, field.Type().Name())
		break
	case reflect.Ptr:
		fieldType := field.Type().Elem()
		value := reflect.New(fieldType)
		if tag.isDefault {
			switch fieldType.Kind() {
			case reflect.String:
				value.Elem().Set(reflect.ValueOf(tag.Default))
				field.Set(value)
				break
			case reflect.Int:
				newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
				value.Elem().Set(reflect.ValueOf(int(newVal)))
				field.Set(value)
				break
			case reflect.Int8:
				newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
				value.Elem().Set(reflect.ValueOf(int8(newVal)))
				field.Set(value)
				break
			case reflect.Int16:
				newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
				value.Elem().Set(reflect.ValueOf(int16(newVal)))
				field.Set(value)
				break
			case reflect.Int32:
				newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
				value.Elem().Set(reflect.ValueOf(int32(newVal)))
				field.Set(value)
				break
			case reflect.Int64:
				newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
				value.Elem().Set(reflect.ValueOf(int(newVal)))
				field.Set(value)
				break
			case reflect.Float32:
				newVal, _ := strconv.ParseFloat(tag.Default, 32)
				value.Elem().Set(reflect.ValueOf(newVal))
				field.Set(value)
				break
			case reflect.Float64:
				newVal, _ := strconv.ParseFloat(tag.Default, 64)
				value.Elem().Set(reflect.ValueOf(newVal))
				field.Set(value)
				break
			case reflect.Bool:
				newVal, _ := strconv.ParseBool(tag.Default)
				value.Elem().Set(reflect.ValueOf(newVal))
				field.Set(value)
				break
			}
		}
		break
	case reflect.Slice, reflect.Array:
		//fmt.Println(tag.Name, fieldKind, field.Type().Name())
		break
	case reflect.Map:
		break
	case reflect.String:
		if tag.isDefault {
			field.Set(reflect.ValueOf(tag.Default))
		}
		break
	case reflect.Int:
		if tag.isDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			field.Set(reflect.ValueOf(int(newVal)))
		}
		break
	case reflect.Int8:
		if tag.isDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			field.Set(reflect.ValueOf(int8(newVal)))
		}
		break
	case reflect.Int16:
		if tag.isDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			field.Set(reflect.ValueOf(int16(newVal)))
		}
		break
	case reflect.Int32:
		if tag.isDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			field.Set(reflect.ValueOf(int32(newVal)))
		}
		break
	case reflect.Int64:
		if tag.isDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			field.Set(reflect.ValueOf(int(newVal)))
		}
		break
	case reflect.Float32:
		if tag.isDefault {
			newVal, _ := strconv.ParseFloat(tag.Default, 32)
			field.Set(reflect.ValueOf(newVal))
		}
		break
	case reflect.Float64:
		if tag.isDefault {
			newVal, _ := strconv.ParseFloat(tag.Default, 64)
			field.Set(reflect.ValueOf(newVal))
		}
		break
	case reflect.Bool:
		if tag.isDefault {
			newVal, _ := strconv.ParseBool(tag.Default)
			field.Set(reflect.ValueOf(newVal))
		}
		break
	}
}
func setDataOn(inputs map[string]interface{}, tag Tags, fieldKind reflect.Kind, field reflect.Value, datesController DatesController) {
	switch fieldKind {
	case reflect.Struct, reflect.Ptr:
		inputType := reflect.TypeOf(inputs[tag.Name])
		inputKind := inputType.Kind()
		fieldType := field.Type().Elem()
		value := reflect.New(fieldType)
		switch inputKind {
		case reflect.Struct, reflect.Ptr:
			rField := setStruct(inputs[tag.Name].(map[string]interface{}), field.Interface(), datesController)
			field.Set(reflect.ValueOf(rField))
		case reflect.String:
			value.Elem().Set(reflect.ValueOf(inputs[tag.Name].(string)))
			field.Set(value)
			break
		case reflect.Int:
			value.Elem().Set(reflect.ValueOf(int(inputs[tag.Name].(int))))
			field.Set(value)
			break
		case reflect.Int8:
			value.Elem().Set(reflect.ValueOf(int8(inputs[tag.Name].(int))))
			field.Set(value)
			break
		case reflect.Int16:
			value.Elem().Set(reflect.ValueOf(int16(inputs[tag.Name].(int))))
			field.Set(value)
			break
		case reflect.Int32:
			value.Elem().Set(reflect.ValueOf(int32(inputs[tag.Name].(int))))
			field.Set(value)
			break
		case reflect.Int64:
			value.Elem().Set(reflect.ValueOf(int64(inputs[tag.Name].(int64))))
			field.Set(value)
			break
		case reflect.Float32:
			value.Elem().Set(reflect.ValueOf(float32(int(inputs[tag.Name].(float32)))))
			field.Set(value)
			break
		case reflect.Float64:
			value.Elem().Set(reflect.ValueOf(int(inputs[tag.Name].(float64))))
			field.Set(value)
			break
		case reflect.Bool:
			value.Elem().Set(reflect.ValueOf(inputs[tag.Name].(bool)))
			field.Set(value)
			break
		}
		break
	case reflect.Slice, reflect.Array:
		if tag.IsObjectID {
			newID, _ := primitive.ObjectIDFromHex(inputs[tag.Name].(string))
			field.Set(reflect.ValueOf(newID))
		} else {
			fmt.Println("slice, array")
		}
		break
	case reflect.Map:
		fmt.Println("map")
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
