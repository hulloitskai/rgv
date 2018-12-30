package stream

import (
	"fmt"
	"net/http"
)

// validateSubreddit ensures that the provided subreddit exists.
func validateSubreddit(name string) (exists bool, err error) {
	req, err := http.NewRequest("GET", "https://www.reddit.com/r/"+name, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "rgv-validator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received non-200 status code (got code %d)",
			res.StatusCode)
	}
	return true, nil
}

// jsonError represents a JSON error.
type jsonError struct {
	Error string `json:"error"`
}
