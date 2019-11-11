package main

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoFields struct {
	CampaignID        string `json:"campaign_id" bson:"campaign_id"`
	TransactionRemark string `json:"transaction_remark" bson:"transaction_remark"`
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

	for i := 0; i < 1000000; i++ {
		oneDoc := MongoFields{
			CampaignID:        "campaign" + strconv.Itoa(i),
			TransactionRemark: strconv.Itoa(i),
		}

		// InsertOne() method Returns mongo.InsertOneResult
		col.InsertOne(context.Background(), oneDoc)
	}

}
