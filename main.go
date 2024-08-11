package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Recipes struct {
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instruction"`
	PublishedAt  time.Time `json:"publishedAt"`
}

func main() {
	r := gin.Default()
	r.Run(":8000")
}
