package common

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func mongoExampleLoad() {
	ctx, canf := context.WithCancel(context.Background())
	defer canf()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println("err:", err)
		return
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("err:", err)
		return
	}
	clo := client.Database("stb")
	// insert(ctx, clo)
	// update(ctx, clo)
	delete(ctx, clo)
	selectOne(ctx, clo)

}

func insert(ctx context.Context, db *mongo.Database) {
	data := bson.M{"room_id": "123", "host_name": "shitingbao", "room_name": "mantou", "room_type": "love", "common": "loss"}
	_, err := db.Collection("chatroom").InsertOne(ctx, data)
	if err != nil {
		log.Println("err:", err)
		return
	}
}

func update(ctx context.Context, db *mongo.Database) {
	_, err := db.Collection("chatroom").UpdateOne(ctx, bson.M{"room_id": "123"}, bson.M{"$set": bson.M{"host_name": "mantou"}})
	if err != nil {
		log.Println("err:", err)
		return
	}

}

func delete(ctx context.Context, db *mongo.Database) {
	if err := db.Collection("chatroom").FindOneAndDelete(ctx, bson.M{"room_id": "123"}).Err(); err != nil {
		log.Println("err:", err)
		return
	}
}

func selectOne(ctx context.Context, db *mongo.Database) {
	var result bson.M
	if err := db.Collection("chatroom").FindOne(ctx, bson.M{"room_id": "123"}).Decode(&result); err != nil {
		log.Println("err:", err)
		return
	}
	log.Println(result)
}
