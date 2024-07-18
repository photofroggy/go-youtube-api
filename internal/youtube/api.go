package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var API_KEY string = ""

func ApiKey() string {
	if len(API_KEY) > 0 {
		return API_KEY
	}

	if err := godotenv.Load(); err != nil {
		return ""
	}

	API_KEY = os.Getenv("YOUTUBE_API_KEY")
	return API_KEY
}

var client *http.Client = &http.Client{Timeout: 10 * time.Second}

type YTApiResponse struct {
	PageInfo struct {
		TotalResults int `json:"totalResults"`
	} `json:"pageInfo"`
	Items []struct {
		Snippet struct {
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
	} `json:"items"`
}

type SearchResult struct {
	Query           string `json:"query"`
	TotalResults    int    `json:"totalResults"`
	ContentCreators string `json:"contentCreators"`
}

func ChannelTitles(data YTApiResponse) string {
	channelTitles := ""
	glue := ""

	for i := range data.Items {
		channelTitle := data.Items[i].Snippet.ChannelTitle

		if strings.Contains(channelTitles, channelTitle) {
			continue
		}

		channelTitles = fmt.Sprintf("%s%s%s", channelTitles, glue, channelTitle)
		glue = ","
	}

	return channelTitles
}

func YoutubeSearch(query string) (SearchResult, error) {
	var result SearchResult

	endpoint := fmt.Sprintf(
		"%s?q=%s&part=snippet&maxResults=20&key=%s",
		"https://www.googleapis.com/youtube/v3/search",
		url.QueryEscape(query),
		ApiKey())

	request, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return result, err
	}

	response, err := client.Do(request)

	if err != nil {
		return result, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	if err != nil {
		return result, err
	}

	var jsonResponse YTApiResponse
	err = json.Unmarshal(body, &jsonResponse)

	if err != nil {
		return result, err
	}

	result.Query = query
	result.TotalResults = jsonResponse.PageInfo.TotalResults
	result.ContentCreators = ChannelTitles(jsonResponse)
	return result, nil
}
