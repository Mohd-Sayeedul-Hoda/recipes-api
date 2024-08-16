package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Recipes struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instruction"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipes

func init() {
	recipes = make([]Recipes, 0)
	file, err := os.ReadFile("recipes.json")
	if err != nil {
		log.Fatal("Error while reading the file ", err)
	}

	err = json.Unmarshal([]byte(file), &recipes)

	if err != nil {
		log.Fatal("Cannot unmarshal the file ", err)
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connect to mongo db server")

	var listOfRecipes []interface{}
	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
}

func NewRecpiesHandler(c *gin.Context) {
	var recipe Recipes

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

func ListRecipiesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipes
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	index := -1
	for i := 0; i < len(recipes); i++ {
		if id == recipes[i].ID {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(recipes); i++ {
		if id == recipes[i].ID {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe Not Found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe Deleted",
	})
}

func SearchRecipeHandler(c *gin.Context) {
	tags := c.Query("tags")
	listOfRecipes := make([]Recipes, 0)

	for i := 0; i < len(recipes); i++ {
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(tags, t) {
				listOfRecipes = append(listOfRecipes, recipes[i])
			}
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	r := gin.Default()
	r.POST("/recipes", NewRecpiesHandler)
	r.GET("/recipes", ListRecipiesHandler)
	r.PUT("/recipes/:id", UpdateRecipeHandler)
	r.DELETE("/recipes/:id", DeleteRecipeHandler)
	r.GET("/recipes/search", SearchRecipeHandler)
	r.Run(":8000")
}
