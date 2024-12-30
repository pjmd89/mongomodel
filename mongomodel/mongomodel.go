package mongomodel

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"time"

	"github.com/pjmd89/goutils/dbutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	dbutils.ModelInterface
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
	o.updateSelf = dbutils.CreateStruct(o.self, true)
}
func (o *Model) SetDBName(dbName string) {
	o.dbName = dbName
}
func (o *Model) Create(inputs map[string]interface{}, opts interface{}) (r interface{}, err error) {
	var createdDate int64 = time.Now().Unix()
	if o.init == false {
		err = errors.New("Not Initialized")
		return r, err
	}
	data, err := SetData(inputs, o.self, DatesController{Created: &createdDate})
	if opts == nil {
		opts = []*options.InsertOneOptions{}
	}
	if err == nil {
		var insertedID any
		insertedID, err = o.conn.Create(data, o.modelName, opts)
		if err == nil && insertedID != nil {
			r, err = o.conn.Read(map[string]interface{}{"_id": insertedID}, o.modelName, nil)
			if err == nil {
				cursor := r.(*mongo.Cursor)
				instance := o.createSliceResult()
				err = cursor.All(context.TODO(), &instance)
				rValue := reflect.ValueOf(instance)
				if rValue.Len() == 1 {
					r = reflect.ValueOf(instance).Index(0).Interface()
				}
			}
		}
		if err == nil && insertedID == nil {
			err = errors.New("Inconsistent data to inserted. Check data sent. No register saved.")
		}
	}

	return r, err
}
func (o *Model) Read(where interface{}, opts interface{}) (r interface{}, err error) {
	var cursor *mongo.Cursor
	if o.init == false {
		err = errors.New("DB not initialized")
		return r, err
	}
	r, err = o.conn.Read(where, o.modelName, opts)
	if err == nil {
		cursor = r.(*mongo.Cursor)
		instance := o.createSliceResult()
		err = cursor.All(context.TODO(), &instance)
		r = instance
	}
	return r, err
}
func (o *Model) RawRead(where interface{}, opts interface{}) (r interface{}, err error) {
	var cursor *mongo.Cursor
	if o.init == false {
		err = errors.New("DB not initialized")
		return r, err
	}
	r, err = o.conn.Read(where, o.modelName, opts)
	if err == nil {
		cursor = r.(*mongo.Cursor)
		instance := make([]map[string]interface{}, 0)
		err = cursor.All(context.TODO(), &instance)
		r = instance
	}
	return r, err
}

func (o *Model) CreateIndex(keys interface{}, indexModelOpts, opts interface{}) (r interface{}, err error) {
	if !o.init {
		err = errors.New("DB not initialized")
		return r, err
	}

	r, err = o.conn.CreateIndex(keys, o.modelName, indexModelOpts, opts)

	return
}

func (o *Model) ListIndexes(opts interface{}) (r interface{}, err error) {
	if !o.init {
		err = errors.New("DB not initialized")
		return r, err
	}

	r, err = o.conn.ListIndexes(o.modelName, opts)
	return
}

func (o *Model) DropIndex(name string, opts interface{}) (r interface{}, err error) {
	if !o.init {
		err = errors.New("DB not initialized")
		return r, err
	}
	r, err = o.conn.DropIndex(name, o.modelName, opts)
	if err == nil {
		re := regexp.MustCompile(`"(\d+)"`)
		r = re.FindStringSubmatch(r.(bson.Raw).String())[1]
	}
	return
}

func (o *Model) Watch(where interface{}, opts interface{}, callback func(any)) (err error) {
	var changeStream *mongo.ChangeStream
	var r any
	if o.init == false {
		err = errors.New("DB not initialized")
		return err
	}
	//pipe := mongo.Pipeline{where}
	r, err = o.conn.Watch(where, o.modelName, opts)
	if err == nil {
		changeStream = r.(*mongo.ChangeStream)
		defer changeStream.Close(context.TODO())

		for changeStream.Next(context.TODO()) {
			callback(changeStream.Current)
		}
	}
	return err
}
func (o *Model) Update(inputs map[string]interface{}, where interface{}, opts interface{}) (r interface{}, err error) {
	var cursor *mongo.Cursor
	var updateDate int64 = time.Now().Unix()
	if o.init == false {
		err = errors.New("DB not initialized")
		return r, err
	}
	data, err := SetData(inputs, o.updateSelf, DatesController{Updated: &updateDate})
	if err == nil {
		r, err = o.conn.Update(Update{Set: data}, where, o.modelName, opts)
		if err == nil {
			cursor = r.(*mongo.Cursor)
			instance := o.createSliceResult()
			cursor.All(context.TODO(), &instance)
			r = instance
		}
	}
	return r, err
}

func (o *Model) FindOneAndUpdate(inputs map[string]interface{}, where interface{}, opts interface{}) (r *mongo.SingleResult, err error) {
	var singleResult any
	if !o.init {
		err = errors.New("DB not initialized")
		return r, err
	}
	updatedDate := time.Now().Unix()
	data, err := SetData(inputs, o.updateSelf, DatesController{Updated: &updatedDate})
	if err == nil {
		singleResult, err = o.conn.FindOneAndUpdate(Update{Set: data}, where, o.modelName, opts)
		if err == nil {
			r = singleResult.(*mongo.SingleResult)
		}
	}
	return r, err
}

func (o *Model) RawUpdate(inputs map[string]interface{}, where interface{}, opts interface{}) (r interface{}, err error) {
	var cursor *mongo.Cursor
	if o.init == false {
		err = errors.New("DB not initialized")
		return r, err
	}

	if err == nil {
		r, err = o.conn.Update(Update{Set: inputs}, where, o.modelName, opts)
		if err == nil {
			cursor = r.(*mongo.Cursor)
			instance := o.createSliceResult()
			cursor.All(context.TODO(), &instance)
			r = instance
		}
	}
	return r, err
}
func (o *Model) Delete(where interface{}, opts interface{}) (r interface{}, err error) {
	var cursor *mongo.Cursor
	if o.init == false {
		err = errors.New("DB not initialized")
		return r, err
	}

	r, err = o.conn.Delete(where, o.modelName, opts)
	if err == nil {
		cursor = r.(*mongo.Cursor)
		instance := o.createSliceResult()
		cursor.All(context.TODO(), &instance)
		r = instance
	}
	return r, err
}
func (o *Model) Replace(inputs map[string]interface{}, where interface{}, opts interface{}) (r interface{}, err error) {
	var cursor *mongo.Cursor
	var updateDate int64 = time.Now().Unix()
	if o.init == false {
		err = errors.New("DB not initialized")
		return r, err
	}

	data, err := SetData(inputs, o.updateSelf, DatesController{Updated: &updateDate})
	if err == nil {
		r, err = o.conn.Replace(data, where, o.modelName, opts)
		if err == nil {
			cursor = r.(*mongo.Cursor)
			instance := o.createSliceResult()
			cursor.All(context.TODO(), &instance)
			r = instance
		}
	}
	return r, err
}
func (o *Model) InterfaceUpdate(inputs interface{}, where interface{}, opts interface{}) (r interface{}, err error) {
	var cursor *mongo.Cursor
	if o.init == false {
		err = errors.New("DB not initialized")
		return r, err
	}

	if err == nil {
		r, err = o.conn.Update(Update{Set: inputs}, where, o.modelName, opts)
		if err == nil {
			cursor = r.(*mongo.Cursor)
			instance := o.createSliceResult()
			cursor.All(context.TODO(), &instance)
			r = instance
		}
	}
	return r, err
}
func (o *Model) InterfaceReplace(data interface{}, where interface{}, opts interface{}) (r interface{}, err error) {
	var cursor *mongo.Cursor
	if o.init == false {
		err = errors.New("DB not initialized")
		return r, err
	}
	if err == nil {
		r, err = o.conn.Replace(data, where, o.modelName, opts)
		if err == nil {
			cursor = r.(*mongo.Cursor)
			instance := o.createSliceResult()
			cursor.All(context.TODO(), &instance)
			r = instance
		}
	}
	return r, err
}
func (o *Model) Repare() (r bool, err error) {
	var result interface{}
	if o.init == false {
		err = errors.New("DB not initialized")
		return
	}
	result, err = o.conn.Read(nil, o.modelName, nil)
	if err != nil {
		return
	}
	cursor := result.(*mongo.Cursor)
	instance := []bson.M{}
	err = cursor.All(context.TODO(), &instance)
	if err != nil {
		return
	}
	err = o.RepareData(o.self, instance)
	return
}
func (o *Model) Count(where interface{}, opts interface{}) (r int64, err error) {
	var count interface{}
	if o.init == false {
		err = errors.New("DB not initialized")
		return r, err
	}
	count, err = o.conn.Count(where, o.modelName, opts)
	if err == nil {
		r = count.(int64)
	}

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
	case reflect.Ptr:
		instance = reflect.TypeOf(o.self).Elem()
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
func (o *Model) RenameCollection(from, to string) {
	conn := o.conn.(*MongoDBConn)
	conn.RenameCollection(from, to)
}
