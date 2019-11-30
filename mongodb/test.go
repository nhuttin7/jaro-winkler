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
	CampaignID        string `json:"campaign_id" bson:"campaign_id"`
	TransactionRemark string `json:"transaction_remark" bson:"transaction_remark"`
}

func SortMapByValue(mKey map[MongoFields]int, m map[MongoFields]float64) []MongoFields {
	type kv struct {
		Key   MongoFields
		Value float64
	}

	var ss []kv
	ssKey := make(map[int]MongoFields, len(mKey))
	for k, v := range m {
		ssKey[mKey[k]] = k
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		if ss[i].Value == ss[j].Value {
			if ss[i].Key.TransactionRemark < ss[j].Key.TransactionRemark {
				return true
			}
		}
		return ss[i].Value > ss[j].Value
	})

	var result []MongoFields
	for _, kv := range ss {
		if len(result) >= 5 {
			break
		}
		result = append(result, kv.Key)
		delete(ssKey, mKey[kv.Key])
		delete(mKey, kv.Key)
	}
	var rest []MongoFields
	var empty MongoFields
	for i := 0; i < len(ss); i++ {
		if ssKey[i] == empty {
			continue
		}
		rest = append(rest, ssKey[i])
	}

	rest = append(result, rest...)
	return rest
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
	col := client.Database("ProcURL").Collection("big")

	filter := bson.M{}
	filter["transaction_remark"] = primitive.Regex{Pattern: "1", Options: "i"}

	op := options.FindOptions{}
	op.SetSort(bson.M{"campaign_id": -1})

	cursor, _ := col.Find(context.Background(), filter, &op)
	bestKey := make(map[MongoFields]int)
	bestResult := make(map[MongoFields]float64)
	t1 := time.Now()
	var key int
	for cursor.Next(context.Background()) {
		var model MongoFields
		cursor.Decode(&model)
		bestResult[model] = JaroWinkler(model.TransactionRemark, "1", 0.7)
		bestKey[model] = key
		key++
	}
	SortMapByValue(bestKey, bestResult)
	fmt.Println("total:", time.Now().Sub(t1).Seconds())
}
