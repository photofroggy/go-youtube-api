package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-youtube-api/internal/youtube"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("q")
	data, err := youtube.YoutubeSearch(query)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err.Error())
		return
	}

	if jsonData, err := json.Marshal(data); err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err.Error())
	} else {
		w.Write(jsonData)
	}
}

func main() {
	http.HandleFunc("/youtubesearch", Handler)
	http.ListenAndServe(":8070", nil)
}
