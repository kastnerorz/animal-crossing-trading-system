package main

import (
	"flag"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/configs"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/routers"
)

func main() {
	flag.StringVar(&configs.MongoURI, "mongo-url", "mongodb://localhost:27017", "")
	flag.StringVar(&configs.MongoCollection, "mongo-collection", "acts-dev", "")
	flag.StringVar(&configs.Port, "port", "8080", "")
	flag.Parse()

	router := routers.SetupRouter()
	router.Run(":" + configs.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
