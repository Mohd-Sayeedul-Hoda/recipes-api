package main

import (
	"context"
	"log"
	"os"
	"time"

	"recipes-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Recipes struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instruction" bson:"instruction"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

var recipesHandler *handlers.RecipesHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	log.Println("Connect to mongo db server")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	log.Println(status)
	recipesHandler = handlers.NewRecipeHandler(ctx, collection, redisClient)
}

func main() {
	r := gin.Default()
	r.POST("/recipes", recipesHandler.NewRecpiesHandler)
	r.GET("/recipes", recipesHandler.ListRecipiesHandler)
	r.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
	r.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	r.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	r.GET("/recipes/search", recipesHandler.SearchRecipeHandler)
	r.Run(":8000")
}
