package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dumpsterfireproject/FoodLogiQ/internal/model"
	"github.com/dumpsterfireproject/FoodLogiQ/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userKey = "userKey"

func main() {
	startServer()
}

func startServer() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	mongoDbName := os.Getenv("MONGO_DB_NAME")
	mongoCollectionName := os.Getenv("MONGO_COLLECTION_NAME")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("error connecting to DB %s", err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("error disconnecting from DB %s", err)
		}
	}()
	eventHandler := service.NewEventHandlerService(service.WithClient(client), service.WithDbName(mongoDbName), service.WithCollectionName(mongoCollectionName))
	authenticationService := service.NewAuthenticationService()

	router := setupHandler(authenticationService, eventHandler)
	router.Run(":8080")
}

func setupHandler(authenticationService service.AuthenticationService, eventHandler service.EventService) *gin.Engine {
	router := gin.Default()
	authenticate := authenticator(authenticationService)
	authenticated := router.Group("/", authenticate)
	authenticated.GET("events", func(c *gin.Context) {
		user := getUser(c)
		events, result := eventHandler.ListEvents(c.Request.Context(), user)
		if result.IsSuccess() {
			c.JSON(result.Status, events)
		} else {
			c.JSON(result.Status, gin.H{"error": result.Err})
		}
	})
	authenticated.GET("events/:id", func(c *gin.Context) {
		user := getUser(c)
		id := c.Param("id")
		event, result := eventHandler.GetEvent(c.Request.Context(), user, id)
		if result.IsSuccess() {
			c.JSON(result.Status, event)
		} else {
			c.JSON(result.Status, gin.H{"error": result.Err})
		}
	})
	authenticated.POST("events", func(c *gin.Context) {
		user := getUser(c)
		event := &model.Event{}
		if err := c.ShouldBindJSON(event); err != nil {
			result := eventHandler.CreateEvent(c.Request.Context(), user, event)
			if result.IsSuccess() {
				c.Status(http.StatusCreated)
			} else {
				c.JSON(result.Status, gin.H{"error": result.Err})
			}
		}
	})
	authenticated.DELETE("events/:id", func(c *gin.Context) {
		user := getUser(c)
		id := c.Param("id")
		result := eventHandler.DeleteEvent(c.Request.Context(), user, id)
		if result.IsSuccess() {
			c.Status(http.StatusOK)
		} else {
			c.JSON(result.Status, gin.H{"error": result.Err})
		}
	})
	return router
}

func authenticator(authenticationService service.AuthenticationService) func(*gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			token = strings.TrimPrefix(token, "Bearer ")
			user := authenticationService.ValidateToken(token)
			if user == nil {
				c.AbortWithStatus(http.StatusForbidden)
			}
			c.Set(userKey, user)
		}
	}
}

func getUser(c *gin.Context) *service.User {
	user, _ := c.Get(userKey)
	switch u := user.(type) {
	case *service.User:
		return u
	default:
		return nil
	}
}
