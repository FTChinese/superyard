package subsapi

import (
	"github.com/FTChinese/superyard/pkg/config"
	"log"
	"net/http"
)

var httpClient = &http.Client{}

type Client struct {
	key     string
	baseURL string
}

func NewClient(debug bool) Client {
	var key string
	var baseURL string

	if debug {
		key = config.MustViperString("web_app.superyard.api_key_dev")
		baseURL = "http://localhost:8200"
	} else {
		key = config.MustViperString("web_app.superyard.api_key_prod")
		baseURL = config.MustViperString("api_url.subscription_v1")
	}

	return Client{
		key:     key,
		baseURL: baseURL,
	}
}

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
