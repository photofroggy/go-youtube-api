# YouTube Search API

This is an example integration of the YouTube API using Go.

# Setup

Make sure `API_KEY` in the `main.go` file is set to a valid API key.

# Running

In a terminal, run the following command:

```sh
$ go run main.go
```

# Usage

You can use this simple app by visiting http://localhost:8080/youtubesearch in your web browser.

The endpoint accepts a search term using the `q` parameters. So, to search for `surfing` you would visit the URL http://localhost:8080?q=surfing

The server returns valid JSON.
