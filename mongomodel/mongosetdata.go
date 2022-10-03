package mongomodel

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pjmd89/goutils/dbutils"
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
			tags := dbutils.GetTags(modelElem.Type().Field(i))
			if tags.IsID {
				field.Set(reflect.ValueOf(id))
				break
			}
		}
	}
	return err
}
func SetData(inputs map[string]interface{}, model interface{}, datesController DatesController) (r interface{}, err error) {
	modelType := reflect.TypeOf(model)
	modelKind := modelType.Kind()

	switch modelKind {
	case reflect.Struct:
		r, err = setStruct(inputs, model, datesController)
	case reflect.Ptr:
		r, err = SetData(inputs, reflect.ValueOf(model).Elem().Interface(), datesController)
	}
	return r, err
}
func setStruct(inputs map[string]interface{}, model interface{}, datesController DatesController) (r interface{}, err error) {
	newModel := reflect.New(reflect.TypeOf(model))

	switch newModel.Elem().Kind() {
	case reflect.Map:
		newModel.Elem().Set(reflect.ValueOf(inputs))
		break
	default:
		for i := 0; i < newModel.Elem().NumField(); i++ {
			field := newModel.Elem().Field(i)
			fieldType := newModel.Elem().Type().Field(i)
			fieldKind := field.Type().Kind()
			tag := dbutils.GetTags(fieldType)
			if inputs == nil {
				err = errors.New("inputs not be nil")
				return
			}
			if tag.CreatedDate && datesController.Created != nil {
				inputs[tag.Name] = *datesController.Created
			}
			if tag.UpdatedDate && datesController.Updated != nil {
				inputs[tag.Name] = *datesController.Updated
			}
			if strings.Trim(tag.Name, " ") != "" {
				if inputs[tag.Name] != nil {
					err = setDataOn(inputs, tag, fieldKind, field, datesController)
				} else {
					err = setNilOn(tag, fieldKind, field, datesController)
				}
				if err != nil {
					newModel = reflect.New(reflect.TypeOf(model))
					break
				}
			}
		}
	}
	return newModel.Interface(), err
}
func setNilOn(tag dbutils.Tags, fieldKind reflect.Kind, field reflect.Value, datesController DatesController) (err error) {
	switch fieldKind {
	case reflect.Struct:
		var rField interface{}
		rField, err = setStruct(map[string]interface{}{}, field.Interface(), datesController)
		field.Set(reflect.ValueOf(rField).Elem())
		break
	case reflect.Ptr:
		fieldType := field.Type().Elem()
		value := reflect.New(fieldType)
		if tag.IsDefault {
			rType := value.Elem().Type()
			switch fieldType.Kind() {
			case reflect.String:
				rValue := reflect.ValueOf(tag.Default)
				value.Elem().Set(rValue.Convert(rType))
				field.Set(value)
				break
			case reflect.Int:
				newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
				rValue := reflect.ValueOf(int(newVal))
				value.Elem().Set(rValue.Convert(rType))
				field.Set(value)
				break
			case reflect.Int8:
				newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
				rValue := reflect.ValueOf(int8(newVal))
				value.Elem().Set(rValue.Convert(rType))
				field.Set(value)
				break
			case reflect.Int16:
				newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
				rValue := reflect.ValueOf(int16(newVal))
				value.Elem().Set(rValue.Convert(rType))
				field.Set(value)
				break
			case reflect.Int32:
				newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
				rValue := reflect.ValueOf(int32(newVal))
				value.Elem().Set(rValue.Convert(rType))
				field.Set(value)
				break
			case reflect.Int64:
				newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
				rValue := reflect.ValueOf(newVal)
				value.Elem().Set(rValue.Convert(rType))
				field.Set(value)
				break
			case reflect.Float32:
				newVal, _ := strconv.ParseFloat(tag.Default, 32)
				rValue := reflect.ValueOf(newVal)
				value.Elem().Set(rValue.Convert(rType))
				field.Set(value)
				break
			case reflect.Float64:
				newVal, _ := strconv.ParseFloat(tag.Default, 64)
				rValue := reflect.ValueOf(newVal)
				value.Elem().Set(rValue.Convert(rType))
				field.Set(value)
				break
			case reflect.Bool:
				newVal, _ := strconv.ParseBool(tag.Default)
				rValue := reflect.ValueOf(newVal)
				value.Elem().Set(rValue.Convert(rType))
				field.Set(value)
				break
			}
		}
		break
	case reflect.Slice, reflect.Array:
		fieldType := field.Type().Elem()
		switch field.Type() {
		case reflect.TypeOf(primitive.ObjectID{}):
			if !tag.IsID {
				field.Set(reflect.ValueOf(primitive.ObjectID{}))
			}
			break
		default:
			field.Set(reflect.MakeSlice(reflect.SliceOf(fieldType), 0, 0))
		}
		break
	case reflect.Map:
		err = fmt.Errorf("attribute \"%s\" Map no set nil", tag.Name)
		break
	case reflect.String:
		if tag.IsDefault {
			field.Set(reflect.ValueOf(tag.Default))
		}
		break
	case reflect.Int:
		if tag.IsDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			field.Set(reflect.ValueOf(int(newVal)))
		}
		break
	case reflect.Int8:
		if tag.IsDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			field.Set(reflect.ValueOf(int8(newVal)))
		}
		break
	case reflect.Int16:
		if tag.IsDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			field.Set(reflect.ValueOf(int16(newVal)))
		}
		break
	case reflect.Int32:
		if tag.IsDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			field.Set(reflect.ValueOf(int32(newVal)))
		}
		break
	case reflect.Int64:
		if tag.IsDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			field.Set(reflect.ValueOf(int(newVal)))
		}
		break
	case reflect.Float32:
		if tag.IsDefault {
			newVal, _ := strconv.ParseFloat(tag.Default, 32)
			field.Set(reflect.ValueOf(newVal))
		}
		break
	case reflect.Float64:
		if tag.IsDefault {
			newVal, _ := strconv.ParseFloat(tag.Default, 64)
			field.Set(reflect.ValueOf(newVal))
		}
		break
	case reflect.Bool:
		if tag.IsDefault {
			newVal, _ := strconv.ParseBool(tag.Default)
			field.Set(reflect.ValueOf(newVal))
		}
		break
	}
	return err
}
func setDataOn(inputs map[string]interface{}, tag dbutils.Tags, fieldKind reflect.Kind, field reflect.Value, datesController DatesController) (err error) {
	switch fieldKind {
	case reflect.Struct:
		var rField interface{}
		rField, err = setStruct(inputs[tag.Name].(map[string]interface{}), field.Interface(), datesController)
		field.Set(reflect.ValueOf(rField).Elem())
		break
	case reflect.Ptr:
		inputType := reflect.TypeOf(inputs[tag.Name])
		inputKind := inputType.Kind()
		fieldType := field.Type().Elem()
		value := reflect.New(fieldType)
		switch inputKind {
		case reflect.Map:
			var rField interface{}
			mField := reflect.New(field.Type().Elem())
			rField, err = setStruct(inputs[tag.Name].(map[string]interface{}), mField.Elem().Interface(), datesController)
			field.Set(reflect.ValueOf(rField))
			break
		case reflect.Array, reflect.Slice:
			if tag.IsObjectID && reflect.TypeOf(inputs[tag.Name]) == reflect.TypeOf(primitive.ObjectID{}) {
				newID := reflect.New(reflect.TypeOf(primitive.ObjectID{}))
				newID.Elem().Set(reflect.ValueOf(inputs[tag.Name]))
				field.Set(newID)
			}
			if !tag.IsObjectID {
				fieldTypex := field.Type().Elem().Elem()

				inputValue := reflect.ValueOf(inputs[tag.Name])
				parseArr := reflect.MakeSlice(reflect.SliceOf(fieldTypex), inputValue.Cap(), inputValue.Cap())
				for i := 0; i < inputValue.Len(); i++ {
					switch field.Type().Elem().Elem().Kind() {
					case reflect.Struct:
						var rField interface{}
						mField := reflect.New(field.Type().Elem().Elem())
						rField, err = setStruct(inputs[tag.Name].([]interface{})[i].(map[string]interface{}), mField.Elem().Interface(), datesController)
						//newArr = reflect.Append(newArr, reflect.ValueOf(rField).Elem())
						parseArr.Index(i).Set(reflect.ValueOf(rField).Elem())
						break
					default:
						parseArr.Index(i).Set(inputValue.Index(i))
					}
				}
				newArr := reflect.New(parseArr.Type())
				newArr.Elem().Set(parseArr)
				field.Set(newArr)
			}
			break
		case reflect.Struct, reflect.Ptr:
			//rField := setStruct(inputs[tag.Name].(map[string]interface{}), field.Interface(), datesController)
			//field.Set(reflect.ValueOf(rField))
		case reflect.String:
			if tag.IsObjectID {
				newID, _ := primitive.ObjectIDFromHex(inputs[tag.Name].(string))
				value.Elem().Set(reflect.ValueOf(newID))
			} else {
				value.Elem().Set(reflect.ValueOf(inputs[tag.Name].(string)))
			}

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
			var newID primitive.ObjectID
			if reflect.TypeOf(inputs[tag.Name]).Kind() == reflect.String {
				newID, _ = primitive.ObjectIDFromHex(inputs[tag.Name].(string))
			} else if reflect.TypeOf(inputs[tag.Name]) == reflect.TypeOf(primitive.ObjectID{}) {
				newID = inputs[tag.Name].(primitive.ObjectID)
			}
			field.Set(reflect.ValueOf(newID))
		} else {
			parseArr := reflect.ValueOf(inputs[tag.Name])
			fieldType := field.Type().Elem()
			newArr := reflect.MakeSlice(reflect.SliceOf(fieldType), 0, 0)
			switch field.Type().Elem().Kind() {
			case reflect.Struct, reflect.Ptr:
				switch reflect.TypeOf(inputs[tag.Name]).Kind() {
				case reflect.Array, reflect.Slice:
					for i := 0; i < parseArr.Len(); i++ {
						var rField interface{}
						mField := reflect.New(field.Type().Elem())
						rField, err = setStruct(inputs[tag.Name].([]interface{})[i].(map[string]interface{}), mField.Elem().Interface(), datesController)
						if err != nil {
							break
						}
						newArr = reflect.Append(newArr, reflect.ValueOf(rField).Elem())
					}
					field.Set(newArr)
				}
				break
			default:

				for i := 0; i < parseArr.Len(); i++ {
					newArr = reflect.Append(newArr, parseArr.Index(i))
				}
				field.Set(newArr)
			}
		}
		break
	case reflect.Map:
		field.Set(reflect.ValueOf(inputs[tag.Name]))
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
	return err
}
