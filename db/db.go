package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// JSON structure containing the URL of the watcher and the list of the seller_currency.
type WatcherURL struct {
	URL            string   `json:"url"`
	SellerCurrency []string `json:"seller_currency"`
}

// JSON structure containing the id of the item.
type ItemID struct {
	Id int `json:"id"`
}

var client *mongo.Client
var db *mongo.Database
var watchersCollection *mongo.Collection
var itemsCollection *mongo.Collection

func InitDB(uri string) error {
	var err error
	client, err = mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	db = client.Database("vinted_go")
	watchersCollection = db.Collection("watchers")
	itemsCollection = db.Collection("items")
	return nil
}

func AppendWatcher(watcher WatcherURL) error {
	_, err := watchersCollection.InsertOne(context.TODO(), watcher)
	return err
}

func ReadWatchers() ([]WatcherURL, error) {
	var watchers []WatcherURL
	cursor, err := watchersCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &watchers); err != nil {
		return nil, err
	}
	return watchers, nil
}

func AppendItemIDs(items []ItemID) error {
	var docs []interface{}
	for _, item := range items {
		docs = append(docs, item)
	}
	_, err := itemsCollection.InsertMany(context.TODO(), docs)
	return err
}

func ItemExists(item ItemID) bool {
	count, err := itemsCollection.CountDocuments(context.TODO(), bson.M{"id": item.Id})
	if err != nil {
		log.Printf("error checking item existence: %v", err)
		return false
	}
	return count > 0
}

func ReadItemIDs() ([]ItemID, error) {
	var items []ItemID
	cursor, err := itemsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &items); err != nil {
		return nil, err
	}
	return items, nil
}
