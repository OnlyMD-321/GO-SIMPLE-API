package storage

import (
	"context"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/OnlyMD-321/GO-SIMPLE-API/internal/models"
)


// MongoHandler handles MongoDB operations
type MongoHandler struct {
	collection *mongo.Collection
}

// // NewMongoHandler initializes the MongoDB handler
// func NewMongoHandler(connectionString string) *MongoHandler {
// 	clientOptions := options.Client().ApplyURI(connectionString)
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to MongoDB: %v", err)
// 	}

// 	collection := client.Database("testdb").Collection("contacts")
// 	return &MongoHandler{collection: collection}
// }

func NewHandler(mongoURI string) (*MongoHandler, error) {
    client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
    if err != nil {
        return nil, err
    }

    err = client.Connect(context.Background())
    if err != nil {
        return nil, err
    }

    db := client.Database("contacts") // Replace with your DB name
	collection := db.Collection("contacts") // Replace with your collection name
	return &MongoHandler{collection: collection}, nil
}

// AddOne adds a document to the collection
func (mh *MongoHandler) AddOne(document interface{}) (*mongo.InsertOneResult, error) {
	return mh.collection.InsertOne(context.TODO(), document)
}

// GetOne retrieves a single document matching the filter
func (mh *MongoHandler) GetOne(result interface{}, filter bson.M) error {
	return mh.collection.FindOne(context.TODO(), filter).Decode(result)
}

// Get retrieves multiple documents matching the filter
func (mh *MongoHandler) Get(filter bson.M) []contact.Contact {
	var contacts []contact.Contact
	cursor, err := mh.collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Failed to fetch contacts: %v", err)
		return nil
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var contact contact.Contact
		if err := cursor.Decode(&contact); err != nil {
			log.Printf("Failed to decode contact: %v", err)
		}
		contacts = append(contacts, contact)
	}

	return contacts
}

// Update updates documents matching the filter
func (mh *MongoHandler) Update(filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	return mh.collection.UpdateOne(context.TODO(), filter, update)
}

// RemoveOne deletes a document matching the filter
func (mh *MongoHandler) RemoveOne(filter bson.M) (*mongo.DeleteResult, error) {
	return mh.collection.DeleteOne(context.TODO(), filter)
}
