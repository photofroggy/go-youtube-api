package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var API_KEY string = "ADD API KEY HERE"

var client *http.Client = &http.Client{Timeout: 10 * time.Second}

type SearchResponse struct {
	PageInfo struct {
		TotalResults int `json:"totalResults"`
	} `json:"pageInfo"`
	Items []struct {
		Snippet struct {
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
	} `json:"items"`
}

func ChannelTitles(data SearchResponse) string {
	channelTitles := ""

	for i := range data.Items {
		glue := ""
		channelTitle := data.Items[i].Snippet.ChannelTitle

		if strings.Contains(channelTitles, channelTitle) {
			continue
		}

		if len(channelTitles) != 0 {
			glue = ","
		}

		channelTitles = fmt.Sprintf("%s%s%s", channelTitles, glue, channelTitle)
	}

	return channelTitles
}

func YoutubeSearch(query string) map[string]any {
	result := make(map[string]any)

	endpoint := fmt.Sprintf(
		"%s?q=%s&part=snippet&maxResults=20&key=%s",
		"https://www.googleapis.com/youtube/v3/search",
		url.QueryEscape(query),
		API_KEY)

	request, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		result["error"] = err.Error()
		return result
	}

	response, err := client.Do(request)

	if err != nil {
		result["error"] = err.Error()
	}

	defer response.Body.Close()

	if err != nil {
		return result
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		result["error"] = err.Error()
		return result
	}

	var jsonResponse SearchResponse
	err = json.Unmarshal(body, &jsonResponse)

	if err != nil {
		result["error"] = err.Error()
		return result
	}

	result["query"] = query
	result["totalResults"] = jsonResponse.PageInfo.TotalResults
	result["contentCreators"] = ChannelTitles(jsonResponse)
	return result
}

func Handler(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("q")
	data := YoutubeSearch(query)
	jsonData, err := json.Marshal(data)

	if err != nil {
		fmt.Fprintf(w, "Encoding error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/youtubesearch", Handler)
	http.ListenAndServe(":8080", nil)
}
