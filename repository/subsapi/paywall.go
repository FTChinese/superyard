package subsapi

import (
	"log"
	"net/http"
)

func (c Client) RefreshPaywall() (*http.Response, error) {
	url := c.baseURL + "/paywall/__refresh"

	log.Printf("Refreshing paywall data at %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.key)

	return httpClient.Do(req)
}
