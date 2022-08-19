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

	var videoDetails []map[string]interface{}

	for _, item := range videoList {
		videoDetails = append(videoDetails, getVideoDetail(apiKey, item))
	}
	scheduledStreams, liveStream, endedStreams, videos := separateVideos(videoDetails)
	fmt.Println("配信予定", scheduledStreams)
	fmt.Println("配信中", liveStream)
	fmt.Println("配信終了", endedStreams)
	fmt.Println("動画", videos)

}

func separateVideos(videoDetails []map[string]interface{}) (scheduledStreams []string, liveStream []string, endedStreams []string, videos []string) {

	for _, item := range videoDetails {
		if videoInfo, ok := item["items"].([]interface{})[0].(map[string]interface{}); ok {

			if _, ok := videoInfo["liveStreamingDetails"]; ok {

				if _, ok := videoInfo["liveStreamingDetails"].(map[string]interface{})["actualEndTime"]; ok {
					if id, ok := videoInfo["id"]; ok {
						endedStreams = append(endedStreams, id.(string))
						// fmt.Println("endedStreams: ", id)
						continue
					}
				}
				if _, ok := videoInfo["liveStreamingDetails"].(map[string]interface{})["actualStartTime"]; ok {
					if id, ok := videoInfo["id"]; ok {
						liveStream = append(liveStream, id.(string))
						// fmt.Println("liveStreams: ", id)
						continue
					}
				}
				if _, ok := videoInfo["liveStreamingDetails"].(map[string]interface{})["scheduledStartTime"]; ok {
					if id, ok := videoInfo["id"]; ok {
						scheduledStreams = append(scheduledStreams, id.(string))
						// fmt.Println("scheduledStreams: ", id)
						continue
					}
				}
			} else {
				if id, ok := videoInfo["id"]; ok {
					videos = append(videos, id.(string))
					// fmt.Println("videos: ", id)
					continue
				}
			}
		}

	}
	return

}

// https://zenn.dev/meihei/articles/1021b1a3f8c226#quota-%E3%81%AE%E7%AF%80%E7%B4%84%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6%EF%BC%88%E3%81%BE%E3%81%A8%E3%82%81%E3%81%A6api%E3%82%92%E5%8F%A9%E3%81%8F%E6%96%B9%E6%B3%95%EF%BC%89
// 引数をvideoIDsにしたものにする
func getVideoDetail(apiKey string, videoId string) (searchData map[string]interface{}) {

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
	// body, _ := ioutil.ReadFile("search.json") //test用
	json.Unmarshal(body, &searchData)
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
	// body, _ := ioutil.ReadFile("search.json") //test用
	json.Unmarshal(body, &searchData)
	return
}
