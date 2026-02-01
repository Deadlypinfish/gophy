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

// GiphySearchResponseFull parses the full Giphy API response with all fields we need
type GiphySearchResponseFull struct {
    Data []struct {
        ID    string `json:"id"`
        Title string `json:"title"`
        Images struct {
            Original struct {
                URL    string `json:"url"`
                Width  string `json:"width"`
                Height string `json:"height"`
            } `json:"original"`
            FixedWidth struct {
                URL string `json:"url"`
            } `json:"fixed_width"`
        } `json:"images"`
    } `json:"data"`
}

// GiphySearchResult is returned to the frontend for display
type GiphySearchResult struct {
    ID         string `json:"id"`
    Title      string `json:"title"`
    OriginalURL string `json:"original_url"`
    PreviewURL  string `json:"preview_url"`
}

// searchGiphy searches Giphy and returns metadata for multiple GIFs
func searchGiphy(q string) []GiphySearchResult {
	giphyAPIKey := os.Getenv("GIPHY_API_KEY")
	if giphyAPIKey == "" {
		panic("GIPHY_API_KEY not set")
	}

	giphyBaseURL := "api.giphy.com/v1/gifs"
	giphySearchURL := giphyBaseURL + "/search"

	encodedQuery := url.QueryEscape(q)
	giphySearchEndpoint := fmt.Sprintf("https://%s?api_key=%s&q=%s&limit=25", giphySearchURL, giphyAPIKey, encodedQuery)

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

	var giphyResp GiphySearchResponseFull
	err = json.Unmarshal(b, &giphyResp)
	if err != nil {
		panic(err)
	}

	// Convert to our result format
	results := make([]GiphySearchResult, 0, len(giphyResp.Data))
	for _, item := range giphyResp.Data {
		results = append(results, GiphySearchResult{
			ID:         item.ID,
			Title:      item.Title,
			OriginalURL: item.Images.Original.URL,
			PreviewURL:  item.Images.FixedWidth.URL,
		})
	}

	return results
}

// downloadGiphyImage downloads a specific GIF from a URL
func downloadGiphyImage(gifURL string) []byte {
	imageReq, _ := http.NewRequest("GET", gifURL, nil)
	imageReq.Header.Set("User-Agent", "Mozilla/5.0")
	client := &http.Client{}
	resp, err := client.Do(imageReq)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return b
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
