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
	searchData := getSearchData(apiKey, "UC7WpJ8eZESNDtO2uALSjigQ", 10)
	videoList := extractVideoList(searchData)

	videoDetails := getVideoDetail(apiKey, videoList)

	scheduledStreams, liveStream, endedStreams, videos := separateVideos(videoDetails)

	fmt.Println("配信予定", scheduledStreams)
	fmt.Println("配信中", liveStream)
	fmt.Println("配信終了", endedStreams)
	fmt.Println("動画", videos)

}

func separateVideos(videoDetails []interface{}) (scheduledStreams []string, liveStream []string, endedStreams []string, videos []string) {
	for _, item := range videoDetails {
		if _, ok := item.(map[string]interface{})["liveStreamingDetails"]; ok {

			if _, ok := item.(map[string]interface{})["liveStreamingDetails"].(map[string]interface{})["actualEndTime"]; ok {
				if id, ok := item.(map[string]interface{})["id"]; ok {
					endedStreams = append(endedStreams, id.(string))
					// fmt.Println("endedStreams: ", id)
					continue
				}
			}
			if _, ok := item.(map[string]interface{})["liveStreamingDetails"].(map[string]interface{})["actualStartTime"]; ok {
				if id, ok := item.(map[string]interface{})["id"]; ok {
					liveStream = append(liveStream, id.(string))
					// fmt.Println("liveStreams: ", id)
					continue
				}
			}
			if _, ok := item.(map[string]interface{})["liveStreamingDetails"].(map[string]interface{})["scheduledStartTime"]; ok {
				if id, ok := item.(map[string]interface{})["id"]; ok {
					scheduledStreams = append(scheduledStreams, id.(string))
					// fmt.Println("scheduledStreams: ", id)
					continue
				}
			}
		} else {
			if id, ok := item.(map[string]interface{})["id"]; ok {
				videos = append(videos, id.(string))
				// fmt.Println("videos: ", id)
				continue
			}
		}
	}

	return

}


func getVideoDetail(apiKey string, videoList []string) (searchData []interface{}) {

	var videoId string

	for i, v := range videoList {
		if (i+1)%50 == 0 {
			videoId += v
			for _, x := range requestVideoDetail(apiKey, videoId) {
				searchData = append(searchData, x)
			}
			videoId = ""
		} else {
			videoId += v + ","
		}
	}

	if len(videoId) != 0 {
		for _, x := range requestVideoDetail(apiKey, videoId) {
			searchData = append(searchData, x)
		}
	}

	return
}

func requestVideoDetail(apiKey string, videoId string) (pieceOfData []interface{}) {
	fmt.Println(videoId)

	q := url.Values{
		"key":  []string{apiKey},
		"part": []string{"liveStreamingDetails"},
		"id":   []string{videoId},
	}
	u := &url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "/youtube/v3/videos",
		RawQuery: q.Encode(),
	}

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal("Error http reques")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Error: status code", resp.StatusCode)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)
	if _, ok := data["items"].([]interface{}); ok {
		pieceOfData = data["items"].([]interface{})
	} else {
		log.Fatal("Error: not included data")
		return
	}
	return
}

func extractVideoList(searchData map[string]interface{}) (videoList []string) {
	for _, item := range searchData["items"].([]interface{}) {
		if video, ok := item.(map[string]interface{})["id"].(map[string]interface{})["videoId"].(string); ok {
			videoList = append(videoList, video)
		}
	}
	return
}

func getSearchData(apiKey string, channelId string, resultCount int) (searchData map[string]interface{}) {

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

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal("Error http reques")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Error: status code", resp.StatusCode)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &searchData)
	return
}
