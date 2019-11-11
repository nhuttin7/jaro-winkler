package main

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoFields struct {
	campaignID        string `json:"campaign_id"`
	transactionRemark string `json:"transaction_remark"`
}

func SortMapByValue(m map[MongoFields]float64) []MongoFields {
	type kv struct {
		Key   MongoFields
		Value float64
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value <= ss[j].Value
	})

	var result []MongoFields
	for _, kv := range ss {
		result = append(result, kv.Key)
	}
	return result
}

func main() {
	// Declare host and port options to pass to the Connect() method
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	fmt.Println("clientOptions type:", reflect.TypeOf(clientOptions), "\n")
	// Connect to the MongoDB and return Client instance
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}

	// Access a MongoDB collection through a database
	col := client.Database("perform_database").Collection("perform")

	filter := bson.M{}
	filter["transaction_remark"] = primitive.Regex{Pattern: "0", Options: "i"}
	t1 := time.Now()

	cursor, _ := col.Find(context.Background(), filter, &options.FindOptions{})
	fmt.Println("find:", time.Now().Sub(t1).Seconds())
	var restModels []MongoFields
	bestResult := make(map[MongoFields]float64)
	for cursor.Next(context.Background()) {
		var model MongoFields
		cursor.Decode(&model)

		if p := JaroWinkler(model.transactionRemark, "O", 0.7); p > 0.7 {
			bestResult[model] = p
		} else {
			restModels = append(restModels, model)
		}
	}
	fmt.Println("after cursor:", time.Now().Sub(t1).Seconds())
	a := SortMapByValue(bestResult)
	for _, i := range a {
		restModels = append([]MongoFields{i}, restModels...)
	}
	fmt.Println("after sorting:", time.Now().Sub(t1).Seconds())
	fmt.Println("total:", time.Now().Sub(t1).Seconds())
	fmt.Println(len(restModels))

}
