package handlers

import (
	"Go-mongo-CRUD/db"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"Go-mongo-CRUD/models"
)

var productCollection *mongo.Collection

var validate = validator.New()

func InitProductHandler() {
	productCollection = db.GetCollection("products")
}

// Create Product
func CreateProduct(c *gin.Context) {
	var newProduct models.Product

	if err := c.BindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	//  Validasi data
	if err := validate.Struct(newProduct); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name":
				if err.Tag() == "required" {
					errors["name"] = "Nama produk wajib diisi"
				} else if err.Tag() == "min" {
					errors["name"] = "Nama produk minimal 2 karakter"
				}
			case "Price":
				if err.Tag() == "required" {
					errors["price"] = "Harga produk wajib diisi"
				} else if err.Tag() == "gt" {
					errors["price"] = "Harga harus lebih dari 0"
				}
			case "Stock":
				if err.Tag() == "required" {
					errors["stock"] = "Stok produk wajib diisi"
				}
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
		return
	}

	newProduct.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := productCollection.InsertOne(ctx, newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan produk"})
		return
	}

	c.JSON(http.StatusCreated, newProduct)
}

// Get All Products
func GetProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := productCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Get Product by ID
func GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var product models.Product
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = productCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Update Product
func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var updateData models.Product
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Validasi data sebelum update
	if err := validate.Struct(updateData); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name":
				if err.Tag() == "required" {
					errors["name"] = "Nama produk wajib diisi"
				} else if err.Tag() == "min" {
					errors["name"] = "Nama produk minimal 2 karakter"
				}
			case "Price":
				if err.Tag() == "required" {
					errors["price"] = "Harga produk wajib diisi"
				} else if err.Tag() == "gt" {
					errors["price"] = "Harga harus lebih dari 0"
				}
			case "Stock":
				if err.Tag() == "required" {
					errors["stock"] = "Stok produk wajib diisi"
				} else if err.Tag() == "gte" {
					errors["stock"] = "Stok tidak boleh negatif"
				}
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":  updateData.Name,
			"price": updateData.Price,
			"stock": updateData.Stock,
		},
	}

	result, err := productCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diupdate"})
}

// Delete Product
func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := productCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk dihapus"})
}
