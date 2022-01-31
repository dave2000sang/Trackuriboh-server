package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"trackuriboh-server/models"

	"go.mongodb.org/mongo-driver/bson"
)

// StartDownload begins downloading card data from TCG API
func StartDownload() error {
	var err error
	// CardSets (called Groups in TCG)
	err = downloadCardSets()
	if err != nil {
		return err
	}

	// TODO: Products

	// TODO: Prices
	return nil
}

func downloadCardSets() error {
	client := &http.Client{}

	// paging params
	limit := 10 // seems like 100 is the limit
	offset := 0

	parsedItems := 0
	totalItems := 2147483647 // arbitrarily large number
	var parsedData []models.CategoryData
	// Paging loop, TODO: spawn new goroutines to speedup
	for parsedItems < totalItems {
		req, err := http.NewRequest("GET", "https://api.tcgplayer.com/catalog/categories/2/groups", nil)
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+AccessToken)
		q := req.URL.Query()
		fmt.Printf("limit: %d, offset: %d\n", limit, offset)
		q.Add("limit", fmt.Sprint(limit))
		q.Add("offset", fmt.Sprint(offset))
		req.URL.RawQuery = q.Encode() // write back to request instance

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		resData, err := io.ReadAll(resp.Body)
		categoryRes := models.CategoryResponse{}
		if err != nil {
			return err
		}
		err = json.Unmarshal(resData, &categoryRes)
		if err != nil {
			return err
		}

		// update DB
		updateDBCategories(parsedData)

		// Continue fetching more responses
		totalItems = categoryRes.TotalItems
		parsedItems += len(categoryRes.Results)
		fmt.Printf("GOT RESPONSE parsed: %d/%d, return count: %d\n", parsedItems, categoryRes.TotalItems, len(categoryRes.Results))
		parsedData = append(parsedData, categoryRes.Results...)
		if totalItems-parsedItems < limit {
			limit = totalItems - parsedItems
		}
		offset = parsedItems
	}
	fmt.Printf("data: %v\n", parsedData)
	fmt.Println(len(parsedData))
	return nil
}

// updateDBCategories updates the database with new data
func updateDBCategories(data []models.CategoryData) error {
	// Parse data into format that matches mongo's document
	categoriesCollection := MongoClient.Database(MONGO_DB).Collection("categories")
	// Write to DB
	doc := bson.D{{"title", "Record of a Shriveled Datum"}, {"text", "No bytes, no problem. Just insert a document, in MongoDB"}}
	result, err := categoriesCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
	// Verify success
	fmt.Println("write result:", result)
	return nil
}
