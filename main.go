package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
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

func main() {
	r := gin.Default()
	r.POST("/recipes", NewRecpiesHandler)
	r.GET("/recipes", ListRecipiesHandler)
	r.PUT("/recipes/:id", UpdateRecipeHandler)
	r.Run(":8000")
}
