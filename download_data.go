package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"trackuriboh-server/models"
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
	req, err := http.NewRequest("GET", "https://api.tcgplayer.com/catalog/categories/2/groups", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+AccessToken)

	// TODO: add paging loop, spawn new goroutines
	// paging params
	limit := 100	// seems like 100 is the limit
	offset := 0
	
	parsedItems := 0
	totalItems := 2147483647	// arbitrarily large number
	for parsedItems < totalItems {
		q := req.URL.Query()
		q.Add("limit", fmt.Sprint(limit))
		q.Add("offset", fmt.Sprint(offset))
		req.URL.RawQuery = q.Encode()		// write back to request instance

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		// var res map[string]interface{}
		// json.NewDecoder(resp.Body).Decode(&res)
		resData, err := io.ReadAll(resp.Body)
		categoryRes := models.CategoryResponse{}
		if err != nil {
			return err
		}
		err = json.Unmarshal(resData, &categoryRes)
		if err != nil {
			return err
		}
		// parse body
		fmt.Println("GOT RESPONSE")
		fmt.Println("totalItems:", categoryRes.TotalItems)
		fmt.Println("results size:", len(categoryRes.Results))
		fmt.Printf("results: %v", categoryRes.Results)
	}
	return nil
}
