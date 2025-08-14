package model

import (
	"fmt"
	"testing"

	"github.com/pjmd89/mongomodel/tests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestAggregate(t *testing.T) {
	//var userResult []tests.TestTypes
	//insertValues := 10
	db := tests.CreateDB()
	testModel := tests.TestTea{}
	testModel.Init(tests.TestTea{}, db)

	total, err := testModel.Count(nil, nil)
	if err != nil {
		t.Fatalf("error al contar los registros: %v", err)
	}

	if total == 0 {
		insertMany(testModel)
	}

	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category"},
			{Key: "average_price", Value: bson.D{{"$avg", "$price"}}},
			{Key: "type_total", Value: bson.D{{"$sum", 1}}},
		}}}

	results, err := testModel.Aggregate(mongo.Pipeline{groupStage}, nil)
	if err != nil {
		t.Fatalf("error during aggregation: %v", err)
	}

	for _, result := range results.([]map[string]any) {
		fmt.Printf("Average price of %v testModel docs: $%v \n", result["_id"], result["average_price"])
		fmt.Printf("Number of %v testModel docs: %v \n\n", result["_id"], result["type_total"])
	}
}

func insertMany(testModel tests.TestTea) {
	docs := []map[string]any{
		{
			"type":     "Masala",
			"category": "black",
			"toppings": []string{"ginger", "pumpkin spice", "cinnamon"},
			"price":    float32(6.75),
		},
		{
			"type":     "Gyokuro",
			"category": "green",
			"toppings": []string{"berries", "milk foam"},
			"price":    float32(5.65),
		},
		{
			"type":     "English Breakfast",
			"category": "black",
			"toppings": []string{"whipped cream", "honey"},
			"price":    float32(5.75),
		},
		{
			"type":     "Sencha",
			"category": "green",
			"toppings": []string{"lemon", "whipped cream"},
			"price":    float32(5.15),
		},
		{
			"type":     "Assam",
			"category": "black",
			"toppings": []string{"milk foam", "honey", "berries"},
			"price":    float32(5.65),
		},
		{
			"type":     "Matcha",
			"category": "green",
			"toppings": []string{"whipped cream", "honey"},
			"price":    float32(6.45),
		},
		{
			"type":     "Earl Grey",
			"category": "black",
			"toppings": []string{"milk foam", "pumpkin spice"},
			"price":    float32(6.15),
		},
		{
			"type":     "Hojicha",
			"category": "green",
			"toppings": []string{"lemon", "ginger", "milk foam"},
			"price":    float32(5.55),
		},
	}

	for _, v := range docs {
		testModel.Create(v, nil)
	}

}
