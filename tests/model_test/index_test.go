package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/pjmd89/mongomodel/tests"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateIndex(t *testing.T) {
	db := tests.CreateDB()
	obj := tests.TestIndex{}
	obj.Init(tests.TestIndex{}, db)

	// dataToInsert := []map[string]interface{}{
	// 	{"theaterId": "7897944646"},
	// }

	indexInput := map[string]any{
		"theaterId": -1,
	}
	opts2 := []*options.CreateIndexesOptions{options.CreateIndexes().SetMaxTime(2 * time.Second)}
	indexModelOpts := options.Index().SetUnique(true)
	_, err := obj.CreateIndex(indexInput, indexModelOpts, opts2)
	if err != nil {
		t.Fatal("error creando índice en la colección: ", err.Error())
	}

	//t.Logf("Se ha creado el índice: %s", createdIndexName.(string))
	// for i := 0; i < len(dataToInsert); i++ {
	// 	r, err := obj.Create(dataToInsert[i], nil)
	// 	if err != nil || r.(tests.TestIndex).TheaterId == "" {
	// 		t.Fatal("Se generó un error al crear el registro: ", err.Error())
	// 	}
	// }

	// where := bson.D{bson.E{Key: "headquarter", Value: bson.M{"$near": bson.A{8.382584154218737, -66.36571055834295}}}}
	// results, err := obj.Read(where, nil)
	// if err != nil || len(results.([]tests.TestIndex)) == 0 {
	// 	t.Fatal("Error en la query: ", err.Error())
	// }
	//t.Log(results.([]tests.TestIndex)[0].Headquarter.Lat, results.([]tests.TestIndex)[0].Headquarter.Lon)
}
func TestListIndex(t *testing.T) {
	db := tests.CreateDB()
	obj := tests.TestIndex{}
	obj.Init(tests.TestIndex{}, db)

	opts := []*options.ListIndexesOptions{options.ListIndexes().SetMaxTime(2 * time.Second)}
	indexes, err := obj.ListIndexes(opts)
	if err != nil {
		t.Fatal("Se generó un error en la query: ", err.Error())
	}

	for i, index := range indexes.([]*mongo.IndexSpecification) {
		fmt.Printf("---- index %d ----\n", i+1)
		fmt.Printf("%#v\n", index)
		// for k, v := range index {
		// 	fmt.Printf("%v: %v\n", k, v)
		// }
		fmt.Println("------------------")
	}

	// where := bson.D{bson.E{Key: "headquarter", Value: bson.M{"$near": bson.A{8.382584154218737, -66.36571055834295}}}}
	// results, err := obj.Read(where, nil)
	// if err != nil || len(results.([]tests.TestIndex)) == 0 {
	// 	t.Fatal("Error en la query: ", err.Error())
	// }
	//t.Log(results.([]tests.TestIndex)[0].Headquarter.Lat, results.([]tests.TestIndex)[0].Headquarter.Lon)
}

func TestDropIndex(t *testing.T) {
	db := tests.CreateDB()
	obj := tests.TestIndex{}
	obj.Init(tests.TestIndex{}, db)
	//TestCreateIndex(t)
	opts := []*options.DropIndexesOptions{options.DropIndexes().SetMaxTime(2 * time.Second)}
	indexBefore, err := obj.DropIndex("theaterId_-1", opts)
	if err != nil {
		t.Fatal("Error al eliminar el índice: ", err.Error())
	}
	fmt.Println(indexBefore)
}
