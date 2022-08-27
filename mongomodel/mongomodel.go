package mongomodel

import (
	"context"
	"errors"
	"reflect"

	"github.com/pjmd89/goutils/dbutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	dbutils.Model
	conn       dbutils.DBInterface
	self       interface{}
	init       bool
	modelName  string
	dbName     string
	updateSelf interface{}
}

func (o *Model) Init(m interface{}, conn dbutils.DBInterface) {
	o.self = m
	o.init = true
	o.conn = conn
	o.modelName = o.getModelName()
	o.updateSelf = setBsonOmitTag(o.self)
}
func (o *Model) SetDBName(dbName string) {
	o.dbName = dbName
}
func (o *Model) GetSkipCollection() []string {
	return o.conn.(*MongoDBConn).SkipCollection
}
func (o *Model) Create(inputs map[string]interface{}, opts interface{}) (r interface{}, err error) {
	if o.init == false {
		err = errors.New("Not Initialized")
		return r, err
	}
	data := SetData(inputs, o.self)
	SetCreatedDate(data)
	r = data
	if opts == nil {
		opts = []*options.InsertOneOptions{}
	}
	insertedID, _ := o.conn.Create(data, o.modelName, opts)
	SetID(r, insertedID.(primitive.ObjectID))
	return r, err
}
func (o *Model) Read(where interface{}, opts interface{}) (r interface{}, err error) {
	var cursor *mongo.Cursor
	if o.init == false {
		err = errors.New("Not Initialized")
		return r, err
	}
	r, err = o.conn.Read(where, o.modelName, opts)
	cursor = r.(*mongo.Cursor)
	instance := o.createSliceResult()
	cursor.All(context.TODO(), &instance)
	r = instance
	return r, err
}
func (o *Model) Update(inputs map[string]interface{}, where interface{}, opts interface{}) (r interface{}, err error) {
	var cursor *mongo.Cursor
	if o.init == false {
		err = errors.New("Not Initialized")
		return r, err
	}
	data := SetData(inputs, o.updateSelf)
	SetUpdatedDate(data)
	r, err = o.conn.Update(data, where, o.modelName, opts)
	cursor = r.(*mongo.Cursor)
	instance := o.createSliceResult()
	cursor.All(context.TODO(), &instance)
	r = instance
	return r, err
}
func (o *Model) Delete(where interface{}, opts interface{}) (r interface{}, err error) {
	var cursor *mongo.Cursor
	if o.init == false {
		err = errors.New("Not Initialized")
		return r, err
	}

	r, err = o.conn.Delete(where, o.modelName, opts)
	cursor = r.(*mongo.Cursor)
	instance := o.createSliceResult()
	cursor.All(context.TODO(), &instance)
	r = instance
	return r, err
}
func (o *Model) Count(where interface{}, opts interface{}) (r int64, err error) {
	var count interface{}
	if o.init == false {
		err = errors.New("Not Initialized")
		return r, err
	}
	count, err = o.conn.Count(where, o.modelName, opts)
	r = count.(int64)
	return
}
func (o *Model) GetModelName() string {
	return o.modelName
}
func (o *Model) createSliceResult() interface{} {
	vType := reflect.TypeOf(o.self).Kind()
	var instance reflect.Type

	switch vType {
	case reflect.Struct:
		instance = reflect.TypeOf(o.self)
		break
	case reflect.Ptr:
		instance = reflect.TypeOf(o.self).Elem()
		break
	}

	i := reflect.MakeSlice(reflect.SliceOf(instance), 0, 0)
	return i.Interface()
}
func (o *Model) getModelName() string {
	vType := reflect.TypeOf(o.self).Kind()
	s := ""
	switch vType {
	case reflect.Struct:
		s = reflect.TypeOf(o.self).Name()
		break
	case reflect.Ptr:
		s = reflect.TypeOf(o.self).Elem().Name()
		break
	}
	return s
}
