package main

import (
	"Go-mongo-CRUD/db"
	"Go-mongo-CRUD/handlers"
	"Go-mongo-CRUD/routes"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	db.MongoConnection()
	defer db.MongoClient.Disconnect(context.TODO())

	fmt.Println("Server running...")

	// Inisialisasi handler setelah koneksi sukses
	handlers.InitProductHandler()

	r := gin.Default()
	routes.ProductRoutes(r)

	r.Run(":8080")
}
