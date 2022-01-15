package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// store API token as global var for now
var AccessToken string

// handleFetch downloads all Yu-gi-oh data from TCGPlayer API
func handleFetch(c *gin.Context) {

}

func setupAuth() (error) {
	// TODO: move these to .env file or smt
	grantType := "client_credentials"
	clientId := "dfe8663b-1fee-41c5-a8be-095a4d4aa765"
	clientSecret := "8f14e005-c8d2-45df-9454-409a4b79b619"

	authURL := "https://api.tcgplayer.com/token"
	
	data := url.Values{
		"grant_type": {grantType},
		"client_id": {clientId},
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

func main() {
    port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// Auth setup
	err := setupAuth();
	if err != nil {
		log.Print(err)
		log.Fatal("Failed to generate access token")
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

	r.GET("/fetch", handleFetch);

	// Listen and serve on defined port
	log.Printf("Listening on port %s", port)
	r.Run(":" + port)
}