package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DroneData struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	DroneID     string             `bson:"droneID" json:"droneID"`
	PublishedAt primitive.DateTime `bson:"publishedAt" json:"publishedAt"`
	Battery     uint16             `bson:"battery" json:"battery"`
	Height      int64              `bson:"height" json:"height"`
	Temperature int64              `bson:"temperature" json:"temperature"`
}

func getPayloadDefault() *DroneData {
	return getPayload(
		primitive.NewObjectID(),
		"Parrot Anafi",
		primitive.NewDateTimeFromTime(time.Now()),
		50,
		60,
		70,
	)
}

func getPayload(id primitive.ObjectID, droneID string, published primitive.DateTime,
	battery uint16, height int64, temperature int64) *DroneData {

	return &DroneData{
		id,
		droneID,
		published,
		battery,
		height,
		temperature,
	}
}

func getClient() (*mongo.Client, context.Context) {
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://thunderpurtz:G2rsb9ae0a64!@cluster0.14i4y.mongodb.net/DronePlatform?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client, ctx
}

// func main() {
// 	clientOptions := options.Client().
// 		ApplyURI("mongodb+srv://thunderpurtz:G2rsb9ae0a64!@cluster0.14i4y.mongodb.net/DronePlatform?retryWrites=true&w=majority")
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	collection := client.Database("DronePlatform").Collection("droneData")

// 	payload := getPayloadDefault()
// 	fmt.Println(payload)

// 	res, err := collection.InsertOne(ctx, payload)

// 	if err != nil {
// 		fmt.Println("error in insert")
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(res.InsertedID)
// 	if err = client.Disconnect(ctx); err != nil {
// 		panic(err)
// 	}
// }
