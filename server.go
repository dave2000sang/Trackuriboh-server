package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global variables
var AccessToken string
var MongoClient *mongo.Client

// TODO: Move these to .env or somewhere safe before productionizing
const MONGO_DB = "trackuriboh"
const MONGO_PASSWORD = "compscilosers"
const MONGO_CLUSTER = "FreeCluster"

// handleDownload downloads all Yu-gi-oh data from TCGPlayer API
func handleDownload(c *gin.Context) {
	StartDownload()
}

func setupApiAuth() error {
	// TODO: move these to .env file or smt
	grantType := "client_credentials"
	clientId := "dfe8663b-1fee-41c5-a8be-095a4d4aa765"
	clientSecret := "8f14e005-c8d2-45df-9454-409a4b79b619"

	authURL := "https://api.tcgplayer.com/token"

	data := url.Values{
		"grant_type":    {grantType},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}

	resp, err := http.PostForm(authURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	AccessToken = fmt.Sprintf("%v", res["access_token"])

	log.Printf("Got access token: %s", AccessToken)
	return nil
}

func setupMongo() error {
	var err error
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://devbitch:%s@freecluster.cwkoe.mongodb.net/%s?retryWrites=true&w=majority", MONGO_PASSWORD, MONGO_CLUSTER))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	MongoClient, err = mongo.Connect(ctx, clientOptions)
	return err
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// TCG API auth setup
	err := setupApiAuth()
	if err != nil {
		log.Print("Failed to generate access token")
		log.Fatal(err)
	}

	if err = setupMongo(); err != nil {
		log.Fatal("Failed to configure mongo cluster")
	}

	// Starts a new Gin instance with no middle-ware
	r := gin.New()

	// Define handlers
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/download", handleDownload)

	// Listen and serve on defined port
	log.Printf("Trackuriboh listening on port %s", port)
	r.Run(":" + port)
}
