package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"math/rand"
)

type GiphyResponse struct {
    Data struct {
        Images struct {
            Original struct {
                URL string `json:"url"`
            } `json:"original"`
        } `json:"images"`
    } `json:"data"`
}

type GiphySearchResponse struct {
    Data []struct {
        Images struct {
            Original struct {
                URL string `json:"url"`
            } `json:"original"`
        } `json:"images"`
    } `json:"data"`
}

func getGiphyImage(q string) []byte {

	giphyAPIKey := os.Getenv("GIPHY_API_KEY")
	if giphyAPIKey == "" {
		panic("GIPHY_API_KEY not set")
	}

	giphyBaseURL := "api.giphy.com/v1/gifs"
	giphySearchURL := giphyBaseURL + "/search"
	// giphyRandomURL := giphyBaseURL + "/random"

	encodedQuery := url.QueryEscape(q)
	giphySearchEndpoint := fmt.Sprintf("https://%s?api_key=%s&q=%s&limit=10", giphySearchURL, giphyAPIKey, encodedQuery)

	req, _ := http.NewRequest("GET", giphySearchEndpoint, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(resp.StatusCode)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var giphyResp GiphySearchResponse
	err = json.Unmarshal(b, &giphyResp)
	if err != nil {
		panic(err)
	}

	randomIndex := rand.Intn(len(giphyResp.Data))
	gifURL := giphyResp.Data[randomIndex].Images.Original.URL

	fmt.Println("GIF URL", gifURL)

	imageReq, _ := http.NewRequest("GET", gifURL, nil)
	imageReq.Header.Set("User-Agent", "Mozilla/5.0")
	client = &http.Client{}
	resp, err = client.Do(imageReq)

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return b
}
