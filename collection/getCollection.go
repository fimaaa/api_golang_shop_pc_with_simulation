package getcollection

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func GetDBName() string {
	return "PCStoreDB"
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(GetDBName()).Collection(collectionName)
	return collection
}
