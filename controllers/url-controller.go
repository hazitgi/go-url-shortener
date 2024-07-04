package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hazi-tgi/go-url-shortner/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type UrlController interface {
	FindAll() ([]*common.URLCollection, error)
	MakeShort(url string) (*common.URLCollection, error)
	FindById(id primitive.ObjectID) (*common.URLCollection, error)
}

type UrlControllerImpl struct {
	dbClient *mongo.Client
	DB       *mongo.Database
}

func NewUrlController(DB *mongo.Client) UrlController {
	return &UrlControllerImpl{
		DB,
		DB.Database("urlShorter"),
	}
}

func (ctrl *UrlControllerImpl) FindAll() ([]*common.URLCollection, error) {
	results := []*common.URLCollection{} // Declare a slice of pointers to URLCollection
	collection := ctrl.DB.Collection("urls")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Printf("Error finding documents: %v", err)
	}
	defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), &results); err != nil {
		fmt.Println("error", err)
	}

	return results, nil
}

// function for saving url
func (ctrl *UrlControllerImpl) MakeShort(url string) (*common.URLCollection, error) {
	newUrlCollection := common.URLCollection{
		Url:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// create a TTL index
	indexOptions := options.Index().SetExpireAfterSeconds(30 * 24 * 60 * 60)
	indexView := ctrl.DB.Collection("urls").Indexes()
	_, err := indexView.CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "createdAt", Value: 1}},
		Options: indexOptions,
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := ctrl.DB.Collection("urls").InsertOne(context.TODO(), newUrlCollection)
	if err != nil {
		fmt.Println("Error inserting document: ", err)
		return nil, err
	}
	newUrlCollection.ID = result.InsertedID.(primitive.ObjectID)
	return &newUrlCollection, nil
}

// function for saving url
func (ctrl *UrlControllerImpl) FindById(id primitive.ObjectID) (*common.URLCollection, error) {
	newUrlCollection := common.NewURLCollection()
	filter := bson.M{"_id": id}
	err := ctrl.DB.Collection("urls").FindOne(context.TODO(), filter).Decode(newUrlCollection)
	if err != nil {
		fmt.Println("Error finding document: ", err)
		return nil, err
	}
	return newUrlCollection, nil
}
