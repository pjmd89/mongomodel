package mongomodel

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
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
		model, err = defineOmitTags(inputs, model)
		if err == nil {
			r, err = setStruct(inputs, model, datesController)
		}
	case reflect.Ptr:
		r, err = SetData(inputs, reflect.ValueOf(model).Elem().Interface(), datesController)
	}
	//getData(r)
	return r, err
}
func defineOmitTags(inputs map[string]interface{}, model interface{}) (r interface{}, err error) {
	fields := []reflect.StructField{}
	newModel := reflect.TypeOf(model)
	for i := 0; i < newModel.NumField(); i++ {
		field := newModel.Field(i)
		tag := dbutils.GetTags(field)
		if inputs == nil {
			err = errors.New("inputs not be nil")
			return
		}
		if inputs[tag.Name] != nil {
			replaceBSON := regexp.MustCompile(`(bson:"[^"]+)(["])`)
			fieldTag := fmt.Sprintf("%v", field.Tag)
			result := replaceBSON.FindString(fieldTag)
			if result != "" {
				replace := regexp.MustCompile(`:`)
				result2 := replace.Split(result, -1)
				tag2 := strings.Replace(result2[1], `"`, "", -1)
				replace2 := regexp.MustCompile(`,`)
				result3 := replace2.Split(tag2, -1)
				var updateTag []string
				for _, sv := range result3 {
					result4 := replace2.Split(sv, -1)
					if result4[0] != "omitempty" {
						updateTag = append(updateTag, sv)
					}
				}
				fieldTag = replaceBSON.ReplaceAllString(fieldTag, `bson:"`+strings.Join(updateTag, ",")+`"`)
				field.Tag = reflect.StructTag(fieldTag)
			}
		}
		fields = append(fields, field)
	}
	newStruct := reflect.StructOf(fields)
	r = reflect.New(newStruct).Elem().Interface()
	return
}
func getData(data interface{}) {
	rType := reflect.TypeOf(data)
	rValue := reflect.ValueOf(data)
	for i := 0; i < rType.Elem().NumField(); i++ {
		field := rValue.Elem().Field(i)
		fieldType := rType.Elem().Field(i)
		if fieldType.Name != "" {
			fmt.Println("field: ", fieldType.Name, field, fieldType.Tag)
		}
	}
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
			//fmt.Println("field: ", fieldType.Name, field, fieldType.Tag)
		}
	}
	//fmt.Println("model: ", newModel.Elem().Interface())
	return newModel.Interface(), err
}
func setNilOn(tag dbutils.Tags, fieldKind reflect.Kind, field reflect.Value, datesController DatesController) (err error) {
	rType := field.Type()
	switch fieldKind {
	case reflect.Struct:
		var rField interface{}
		rField, err = setStruct(map[string]interface{}{}, field.Interface(), datesController)
		field.Set(reflect.ValueOf(rField).Elem())
		//fmt.Println("rfield: ", reflect.ValueOf(rField).Interface(), "field: ", field.Interface())
		break
	case reflect.Ptr:
		fieldType := field.Type().Elem()
		value := reflect.New(fieldType)
		if tag.IsDefault {
			rType = value.Elem().Type()
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
			rValue := reflect.ValueOf(tag.Default)
			field.Set(rValue.Convert(rType))
		}
		break
	case reflect.Int:
		if tag.IsDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			rValue := reflect.ValueOf(int(newVal))
			field.Set(rValue.Convert(rType))
		}
		break
	case reflect.Int8:
		if tag.IsDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			rValue := reflect.ValueOf(int8(newVal))
			field.Set(rValue.Convert(rType))
		}
		break
	case reflect.Int16:
		if tag.IsDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			rValue := reflect.ValueOf(int16(newVal))
			field.Set(rValue.Convert(rType))
		}
		break
	case reflect.Int32:
		if tag.IsDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			rValue := reflect.ValueOf(int32(newVal))
			field.Set(rValue.Convert(rType))
		}
		break
	case reflect.Int64:
		if tag.IsDefault {
			newVal, _ := strconv.ParseInt(tag.Default, 10, 64)
			rValue := reflect.ValueOf(int64(newVal))
			field.Set(rValue.Convert(rType))
		}
		break
	case reflect.Float32:
		if tag.IsDefault {
			newVal, _ := strconv.ParseFloat(tag.Default, 32)
			rValue := reflect.ValueOf(newVal)
			field.Set(rValue.Convert(rType))
		}
		break
	case reflect.Float64:
		if tag.IsDefault {
			newVal, _ := strconv.ParseFloat(tag.Default, 64)
			rValue := reflect.ValueOf(newVal)
			field.Set(rValue.Convert(rType))
		}
		break
	case reflect.Bool:
		if tag.IsDefault {
			newVal, _ := strconv.ParseBool(tag.Default)
			rValue := reflect.ValueOf(newVal)
			field.Set(rValue.Convert(rType))
		}
		break
	}
	return err
}
func setDataOn(inputs map[string]interface{}, tag dbutils.Tags, fieldKind reflect.Kind, field reflect.Value, datesController DatesController) (err error) {
	rType := field.Type()
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
				rValue := reflect.ValueOf(inputs[tag.Name].(string))
				value.Elem().Set(reflect.ValueOf(rValue.Convert(rType)))
			}

			field.Set(value)
			break
		case reflect.Int:
			rValue := reflect.ValueOf(int(inputs[tag.Name].(int)))
			value.Elem().Set(reflect.ValueOf(rValue.Convert(rType)))
			field.Set(value)
			break
		case reflect.Int8:
			rValue := reflect.ValueOf(int8(inputs[tag.Name].(int8)))
			value.Elem().Set(reflect.ValueOf(rValue.Convert(rType)))
			field.Set(value)
			break
		case reflect.Int16:
			rValue := reflect.ValueOf(int16(inputs[tag.Name].(int16)))
			value.Elem().Set(reflect.ValueOf(rValue.Convert(rType)))
			field.Set(value)
			break
		case reflect.Int32:
			rValue := reflect.ValueOf(int32(inputs[tag.Name].(int32)))
			value.Elem().Set(reflect.ValueOf(rValue.Convert(rType)))
			field.Set(value)
			break
		case reflect.Int64:
			rValue := reflect.ValueOf(int64(inputs[tag.Name].(int64)))
			value.Elem().Set(reflect.ValueOf(rValue.Convert(rType)))
			field.Set(value)
			break
		case reflect.Float32:
			rValue := reflect.ValueOf(float32(int(inputs[tag.Name].(float32))))
			value.Elem().Set(reflect.ValueOf(rValue.Convert(rType)))
			field.Set(value)
			break
		case reflect.Float64:
			rValue := reflect.ValueOf(float64(int(inputs[tag.Name].(float64))))
			value.Elem().Set(reflect.ValueOf(rValue.Convert(rType)))
			field.Set(value)
			break
		case reflect.Bool:
			value.Elem().Set(reflect.ValueOf(inputs[tag.Name].(bool)))
			field.Set(value)
			break
		}
		break
	case reflect.Slice, reflect.Array:
		if tag.IsObjectID && field.Type() == reflect.TypeOf(primitive.ObjectID{}) {
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
				if reflect.TypeOf(primitive.ObjectID{}) == fieldType && reflect.TypeOf(inputs[tag.Name]).Elem().Kind() == reflect.Interface {
					for i := 0; i < parseArr.Len(); i++ {
						//fmt.Println(parseArr.Index(i))
						newArr = reflect.Append(newArr, reflect.ValueOf(parseArr.Index(i).Interface().(primitive.ObjectID)))
					}
				} else {
					for i := 0; i < parseArr.Len(); i++ {
						newArr = reflect.Append(newArr, parseArr.Index(i))
					}
				}
				field.Set(newArr)
			}
		}
		break
	case reflect.Map:
		field.Set(reflect.ValueOf(inputs[tag.Name]))
		break
	case reflect.String:
		rValue := reflect.ValueOf(inputs[tag.Name].(string))
		field.Set(rValue.Convert(rType))
		break
	case reflect.Int:
		rValue := reflect.ValueOf(inputs[tag.Name].(int))
		field.Set(rValue.Convert(rType))
		break
	case reflect.Int8:
		rValue := reflect.ValueOf(inputs[tag.Name].(int8))
		field.Set(rValue.Convert(rType))
		break
	case reflect.Int16:
		rValue := reflect.ValueOf(inputs[tag.Name].(int16))
		field.Set(rValue.Convert(rType))
		break
	case reflect.Int32:
		rValue := reflect.ValueOf(inputs[tag.Name].(int32))
		field.Set(rValue.Convert(rType))
		break
	case reflect.Int64:
		rValue := reflect.ValueOf(inputs[tag.Name].(int64))
		field.Set(rValue.Convert(rType))
		break
	case reflect.Float32:
		rValue := reflect.ValueOf(inputs[tag.Name].(float32))
		field.Set(rValue.Convert(rType))
		break
	case reflect.Float64:
		rValue := reflect.ValueOf(inputs[tag.Name].(float64))
		field.Set(rValue.Convert(rType))
		break
	case reflect.Bool:
		rValue := reflect.ValueOf(inputs[tag.Name].(bool))
		field.Set(rValue.Convert(rType))
		break
	}
	return err
}
