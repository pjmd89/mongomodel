package mongomodel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/pjmd89/goutils/dbutils"
	"github.com/pjmd89/goutils/jsonutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewConnWithURIStruct(config URIData) (r dbutils.DBInterface) {
	db := &MongoDBConn{
		DB: dbutils.DB{
			Host:     config.Host,
			Port:     config.Port,
			User:     config.User,
			Pass:     config.Pass,
			DataBase: config.DataBase,
		},
	}
	db.database = db.DataBase
	r = db
	r.Connect()
	return r
}

func NewConn(configPath *string) (r dbutils.DBInterface) {
	db := &MongoDBConn{}
	configFile := "./etc/db/db.json"

	if configPath != nil {
		configFile = *configPath
	}
	jsonutils.GetJson(configFile, db)
	db.database = db.DataBase
	r = db
	r.Connect()

	return r
}

func (o *MongoDBConn) Connect() (err error) {
	o.tryingCounter = 0
	uri := o.getURI()
	sleep := 1
	monitor := &event.PoolMonitor{
		Event: o.monitor,
	}

	opts := options.Client().ApplyURI(uri)
	opts.SetPoolMonitor(monitor).SetHeartbeatInterval(time.Duration(sleep) * time.Second)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(sleep)*time.Second)
	conn, _ := mongo.Connect(ctx, opts)
	defer cancel()
	o.client = conn

	err = o.ping()
	if err != nil {
		fmt.Println("error connecting to MongoDB")
		o.client.Disconnect(context.TODO())
	} else {
		fmt.Println("connected to MongoDB")
	}

	return err
}
func (o *MongoDBConn) Close() error {
	return o.client.Disconnect(context.TODO())
}
func (o *MongoDBConn) GetClient() interface{} {
	return o.client
}
func (o *MongoDBConn) RenameCollection(from, to string) {

	rename := bson.D{
		{Key: "renameCollection", Value: from},
		{Key: "to", Value: to},
	}
	var result bson.M
	err := o.client.Database("admin").RunCommand(context.TODO(), rename).Decode(&result)
	if err != nil {
		log.Println("Rename collection error: " + err.Error())
	}

}
func (o *MongoDBConn) monitor(evt *event.PoolEvent) {
	if o.Reconnect {
		sleep := 5
		switch evt.Type {
		case event.PoolClosedEvent, event.ConnectionClosed:
			if o.tryingCounter == 0 {
				o.tryingCounter++
				fmt.Println("Trying again in ", sleep, " seconds")
				time.Sleep(time.Duration(sleep) * time.Second)
				o.Connect()
			}
		}
	}
}
func (o *MongoDBConn) getURI() string {
	uri := "mongodb://"
	credentials := o.User + ":" + o.Pass + "@"
	instance := o.Host + ":" + o.Port

	if strings.Trim(o.Pass, " ") == "" {
		credentials = o.User + "@"
	}
	if strings.Trim(o.User, " ") == "" {
		credentials = ""
	}
	if strings.Trim(o.Port, " ") == "" {
		instance = o.Host
	}
	uri += credentials + instance

	return uri + "/" + o.DataBase + "?maxPoolSize=20&w=majority"
}
func (o *MongoDBConn) ping() error {
	err := o.client.Ping(context.TODO(), nil)
	return err
}
func (o *MongoDBConn) Create(inputs interface{}, collection string, opts interface{}) (results interface{}, err error) {
	checkCollection, database, collection := o.CheckCollection(collection)
	checkOpts := true
	if !checkCollection {
		err = errors.New("no collection specified")
		return nil, err
	}
	if opts == nil {
		opts = []*options.InsertOneOptions{}
	}
	optsKind := reflect.ValueOf(opts).Kind()

	switch optsKind {
	case reflect.Slice:
		for i, v := range opts.([]*options.InsertOneOptions) {
			optsType := reflect.ValueOf(v).Type()
			if optsType != reflect.TypeOf(&options.InsertOneOptions{}) {
				err = fmt.Errorf("opts %d value is not *options.InsertOneOptions", i)
				checkOpts = false
				break
			}
		}
		break
	default:
		err = errors.New("opts is not a Slice")
		checkOpts = false
	}

	if checkOpts {
		coll := o.client.Database(database).Collection(collection)
		var r *mongo.InsertOneResult
		r, err = coll.InsertOne(context.TODO(), inputs, opts.([]*options.InsertOneOptions)...)
		if err == nil {
			results = r.InsertedID
		}
	}
	return results, err
}
func (o *MongoDBConn) Read(where interface{}, collection string, opts interface{}) (results interface{}, err error) {
	var cursor *mongo.Cursor
	results = cursor
	checkCollection, database, collection := o.CheckCollection(collection)
	if !checkCollection {
		err = errors.New("no collection specified")
		return nil, err
	}
	coll := o.client.Database(database).Collection(collection)
	if where == nil {
		where = bson.M{}
	}

	if opts != nil {
		if err = o.evaluateType(opts, []*options.FindOptions{}); err != nil {
			return
		}
	} else {
		opts = []*options.FindOptions{}
	}

	cursor, err = coll.Find(context.TODO(), where, opts.([]*options.FindOptions)...)
	if err == nil {
		results = cursor
	}
	return results, err
}

func (p *MongoDBConn) evaluateType(typpe, expected interface{}) (err error) {
	inputValueType := reflect.ValueOf(typpe).Type()
	expectedValueType := reflect.TypeOf(expected)
	if inputValueType != expectedValueType {
		err = fmt.Errorf("opts type must be %s not %s", expectedValueType.String(), inputValueType.String())
	}
	return
}

func (o *MongoDBConn) Aggregate(pipeline any, collection string, opts any) (results any, aggregationError error) {
	checkCollection, database, collection := o.CheckCollection(collection)
	if !checkCollection {
		aggregationError = errors.New("no collection specified")
		return nil, aggregationError
	}

	coll := o.client.Database(database).Collection(collection)

	if opts != nil {
		if aggregationError = o.evaluateType(opts, []*options.AggregateOptions{}); aggregationError != nil {
			return
		}
	} else {
		opts = []*options.AggregateOptions{}
	}

	if pipeline == nil {
		pipeline = mongo.Pipeline{}
	}

	return coll.Aggregate(context.TODO(), pipeline, opts.([]*options.AggregateOptions)...)
}

func (o *MongoDBConn) Watch(where interface{}, collection string, opts interface{}) (results interface{}, err error) {
	var cursor *mongo.ChangeStream
	results = cursor
	checkCollection, database, collection := o.CheckCollection(collection)
	checkOpts := true
	if !checkCollection {
		err = errors.New("no collection specified")
		return nil, err
	}
	coll := o.client.Database(database).Collection(collection)
	if where == nil {
		where = bson.M{}
	}

	if opts == nil {
		opts = []*options.ChangeStreamOptions{}
	}
	optsKind := reflect.ValueOf(opts).Kind()

	switch optsKind {
	case reflect.Slice:
		for i, v := range opts.([]*options.ChangeStreamOptions) {
			optsType := reflect.ValueOf(v).Type()
			if optsType != reflect.TypeOf(&options.ChangeStreamOptions{}) {
				err = fmt.Errorf("opts %d value is not *options.FindOptions", i)
				checkOpts = false
				break
			}
		}
		break
	default:
		err = errors.New("opts is not a Slice")
		checkOpts = false
	}
	if checkOpts {
		cursor, err = coll.Watch(context.TODO(), where, opts.([]*options.ChangeStreamOptions)...)
		if err == nil {
			results = cursor
		}
	}
	return results, err
}

func (o *MongoDBConn) Update(inputs interface{}, where interface{}, collection string, opts interface{}) (results interface{}, err error) {
	var cursor *mongo.Cursor
	checkCollection, database, collection := o.CheckCollection(collection)
	if !checkCollection {
		err = errors.New("no collection specified")
		return nil, err
	}

	collOpts := options.Collection().SetBSONOptions(
		&options.BSONOptions{
			OmitZeroStruct: true,
		},
	)

	coll := o.client.Database(database).Collection(collection, collOpts)
	if where == nil {
		where = bson.M{}
	}
	if opts != nil {
		if err = o.evaluateType(opts, []*options.UpdateOptions{}); err != nil {
			return
		}
	} else {
		opts = []*options.UpdateOptions{}
	}
	_, err = coll.UpdateOne(context.TODO(), where, inputs, opts.([]*options.UpdateOptions)...)
	cursor, _ = coll.Find(context.TODO(), where)
	results = cursor
	if err != nil {
		var x *mongo.Cursor
		results = x
	}

	return results, err
}

func (o *MongoDBConn) Replace(inputs interface{}, where interface{}, collection string, opts interface{}) (results interface{}, err error) {
	var cursor *mongo.Cursor
	checkCollection, database, collection := o.CheckCollection(collection)
	checkOpts := true
	if !checkCollection {
		err = errors.New("no collection specified")
		return nil, err
	}
	coll := o.client.Database(database).Collection(collection)
	if where == nil {
		where = bson.M{}
	}

	if opts == nil {
		opts = []*options.ReplaceOptions{}
	}
	optsKind := reflect.ValueOf(opts).Kind()

	switch optsKind {
	case reflect.Slice:
		for i, v := range opts.([]*options.ReplaceOptions) {
			optsType := reflect.ValueOf(v).Type()
			if optsType != reflect.TypeOf(&options.ReplaceOptions{}) {
				err = fmt.Errorf("opts %d value is not *options.ReplaceOptions", i)
				checkOpts = false
				break
			}
		}
		break
	default:
		err = errors.New("opts is not a Slice")
		checkOpts = false
	}
	if checkOpts {
		_, err = coll.ReplaceOne(context.TODO(), where, inputs, opts.([]*options.ReplaceOptions)...)
		cursor, _ = coll.Find(context.TODO(), where)
		results = cursor
		if err != nil {
			var x *mongo.Cursor
			results = x
		}
	}
	return results, err
}
func (o *MongoDBConn) Delete(where interface{}, collection string, opts interface{}) (results interface{}, err error) {
	var cursor *mongo.Cursor
	checkCollection, database, collection := o.CheckCollection(collection)
	checkOpts := true
	if !checkCollection {
		err = errors.New("no collection specified")
		return nil, err
	}
	coll := o.client.Database(database).Collection(collection)
	if where == nil {
		where = bson.M{}
	}

	if opts == nil {
		opts = []*options.DeleteOptions{}
	}
	optsKind := reflect.ValueOf(opts).Kind()

	switch optsKind {
	case reflect.Slice:
		for i, v := range opts.([]*options.DeleteOptions) {
			optsType := reflect.ValueOf(v).Type()
			if optsType != reflect.TypeOf(&options.DeleteOptions{}) {
				err = fmt.Errorf("opts %d value is not *options.DeleteOptions", i)
				checkOpts = false
				break
			}
		}
		break
	default:
		err = errors.New("opts is not a Slice")
		checkOpts = false
	}
	if checkOpts {
		cursor, _ = coll.Find(context.TODO(), where)
		_, err = coll.DeleteMany(context.TODO(), where, opts.([]*options.DeleteOptions)...)
		results = cursor
		if err != nil {
			var x *mongo.Cursor
			results = x
		}
	}
	return results, err
}
func (o *MongoDBConn) Count(where interface{}, collection string, opts interface{}) (results interface{}, err error) {
	var count int64
	checkCollection, database, collection := o.CheckCollection(collection)
	checkOpts := true
	if !checkCollection {
		err = errors.New("no collection specified")
		return nil, err
	}
	coll := o.client.Database(database).Collection(collection)
	if where == nil {
		where = bson.M{}
	}

	if opts == nil {
		opts = []*options.CountOptions{}
	}
	optsKind := reflect.ValueOf(opts).Kind()

	switch optsKind {
	case reflect.Slice:
		for i, v := range opts.([]*options.CountOptions) {
			optsType := reflect.ValueOf(v).Type()
			if optsType != reflect.TypeOf(&options.CountOptions{}) {
				err = fmt.Errorf("opts %d value is not *options.CountOptions", i)
				checkOpts = false
				break
			}
		}
		break
	default:
		err = errors.New("opts is not a Slice")
		checkOpts = false
	}
	if checkOpts {
		count, err = coll.CountDocuments(context.TODO(), where, opts.([]*options.CountOptions)...)
		results = count
	}
	return results, err
}
func (o *MongoDBConn) CheckCollection(currentCollection string) (r bool, database string, collection string) {
	database = strings.Trim(o.GetDatabase(), " ")
	collection = strings.Trim(currentCollection, " ")
	if o.OnDatabase != nil {
		newDB := strings.Trim(o.OnDatabase(database, collection), " ")
		if newDB != "" {
			database = newDB
		}
	}
	if collection != "" {
		r = true
	}
	return r, database, collection
}
func (o *MongoDBConn) SetDatabase(database string) {
	o.database = database
}
func (o *MongoDBConn) SetCollection(collection string) {
	o.collection = collection
}
func (o *MongoDBConn) GetDatabase() (r string) {
	return o.database
}
