package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	searchData := getsearchData(apiKey, "UC7WpJ8eZESNDtO2uALSjigQ", 2)
	videoList := getVideoList(searchData)
	for _, item := range videoList {
		fmt.Println(item)
	}
	// fmt.Println(searchData)
	// fmt.Println()
	// fmt.Println(searchData["items"])

}

func getVideoList(searchData map[string]interface{}) (videoList []string) {
	for _, item := range searchData["items"].([]interface{}) {
		if video, ok := item.(map[string]interface{})["id"].(map[string]interface{})["videoId"].(string); ok {
			videoList = append(videoList, video)
		}
	}
	return
}

func getsearchData(apiKey string, channelId string, resultCount int) (searchData map[string]interface{}) {

	q := url.Values{
		"key":        []string{apiKey},
		"part":       []string{"id"},
		"channelId":  []string{channelId},
		"order":      []string{"date"},
		"maxResults": []string{strconv.Itoa(resultCount)},
	}
	u := &url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "youtube/v3/search",
		RawQuery: q.Encode(),
	}

	fmt.Println(u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal("Error http reques")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Error: status code", resp.StatusCode)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	// body, _ := ioutil.ReadFile("search.json") //test用
	json.Unmarshal(body, &searchData)
	return
}
