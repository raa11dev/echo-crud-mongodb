package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/raa11dev/crud-echo/database"
	"github.com/raa11dev/crud-echo/models"
	"github.com/raa11dev/crud-echo/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection = database.GetCollection(database.DB, "products")
var validate = validator.New()

func CreateProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var product models.Product
	defer cancel()

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ProductResponses{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	if validationErr := validate.Struct(&product); validationErr != nil {
		return c.JSON(http.StatusInternalServerError, responses.ProductResponses{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": validationErr.Error()},
		})
	}

	var newStatus string
	if product.Status == true {
		newStatus = "ACTIVE"
	} else {
		newStatus = "DEACTIVE"
	}
	newProduct := models.ProductCreate{
		Id:           product.Id,
		Product_name: product.Product_name,
		Status:       newStatus,
	}

	result, err := productCollection.InsertOne(ctx, newProduct)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ProductResponses{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, responses.ProductResponses{
		Status:  http.StatusOK,
		Message: "Product created successfully",
		Data:    &echo.Map{"data": result},
	})
}

func UpdateProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	productId := c.Param("id")
	var product models.Product
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(productId)

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ProductResponses{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	if validationErr := validate.Struct(&product); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.ProductResponses{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": validationErr.Error()},
		})
	}

	var newStatus string
	if product.Status == true {
		newStatus = "ACTIVE"
	} else {
		newStatus = "DEACTIVE"
	}

	update := bson.M{"id": product.Id, "product_name": product.Product_name, "status": newStatus}

	result, err := productCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ProductResponses{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	var updatedUser models.ProductCreate
	if result.MatchedCount == 1 {
		err := productCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.ProductResponses{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &echo.Map{"data": err.Error()},
			})
		}
	}

	return c.JSON(http.StatusOK, responses.ProductResponses{
		Status:  http.StatusOK,
		Message: "Product edited successfully",
		Data:    &echo.Map{"data": updatedUser},
	})
}

func DeleteProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	productId := c.Param("id")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(productId)

	result, err := productCollection.DeleteOne(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ProductResponses{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	if result.DeletedCount < 1 {
		return c.JSON(http.StatusNotFound, responses.ProductResponses{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    &echo.Map{"data": "Product not found"},
		})
	}

	return c.JSON(http.StatusOK, responses.ProductResponses{
		Status:  http.StatusOK,
		Message: "Product deleted successfully",
		Data:    &echo.Map{"data": result},
	})
}

func GetProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	productId := c.Param("id")
	var user models.ProductCreate
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(productId)

	err := productCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ProductResponses{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, responses.ProductResponses{
		Status:  http.StatusOK,
		Message: "Product found",
		Data:    &echo.Map{"data": user},
	})
}
