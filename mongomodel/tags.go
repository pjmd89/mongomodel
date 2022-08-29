package mongomodel

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Tags struct {
	Name        string
	Default     string
	isDefault   bool
	IsID        bool
	IsObjectID  bool
	CreatedDate bool
	UpdatedDate bool
}

func setBsonOmitTag(instance interface{}) (r interface{}) {
	valueOf := reflect.ValueOf(instance)
	typeOf := valueOf.Type()
	structFields := make([]reflect.StructField, 0)
	for i := 0; i < typeOf.NumField(); i++ {
		tag := fmt.Sprintf("%v", typeOf.Field(i).Tag)
		tagFind := regexp.MustCompile(`bson:"[^"\-]+"`)
		notFind := regexp.MustCompile(`omitempty`)
		result := tagFind.FindString(tag)
		if !notFind.MatchString(result) && strings.Trim(result, " ") != "" {
			replace := regexp.MustCompile(`(bson:"[^"]+)(["])`)
			tag = replace.ReplaceAllString(tag, "$1,omitempty\"")
		}
		structFields = append(structFields, typeOf.Field(i))
		structFields[i].Tag = reflect.StructTag(tag)
	}
	newType := reflect.StructOf(structFields)
	newStruct := valueOf.Convert(newType).Interface()
	r = newStruct
	return r
}
func getTags(field reflect.StructField) (r Tags) {
	tag := field.Tag.Get("gql")
	if tag != "" {
		tagSplit := strings.Split(tag, ",")
		for _, v := range tagSplit {
			dataSplit := strings.Split(v, "=")
			switch dataSplit[0] {
			case "name":
				r.Name = dataSplit[1]
				break
			case "default":
				r.Default = dataSplit[1]
				r.isDefault = true
				break
			case "id":
				isTrue := false
				if strings.Trim(dataSplit[1], " ") == "true" {
					isTrue = true
				}
				r.IsID = isTrue
				break
			case "objectID":
				isTrue := false
				if strings.Trim(dataSplit[1], " ") == "true" {
					isTrue = true
				}
				r.IsObjectID = isTrue
				break
			case "createdDate":
				isTrue := false
				if strings.Trim(dataSplit[1], " ") == "true" {
					isTrue = true
				}
				r.CreatedDate = isTrue
				break
			case "updatedDate":
				isTrue := false
				if strings.Trim(dataSplit[1], " ") == "true" {
					isTrue = true
				}
				r.UpdatedDate = isTrue
				break
			}
		}
	}
	return r
}
