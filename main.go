package main

import (
	"net/http"
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
}

func NewRecpiesHandler(c *gin.Context) {
	var recipe Recipes

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

func main() {
	r := gin.Default()
	r.POST("/recipes", NewRecpiesHandler)
	r.Run(":8000")
}
