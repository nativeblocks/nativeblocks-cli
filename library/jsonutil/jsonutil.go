package jsonutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchJSONFromURL(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch JSON from %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP response status: %s from %s", resp.Status, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON from %s: %v", url, err)
	}
	return nil
}
