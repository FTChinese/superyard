package subsapi

import (
	"fmt"
	"github.com/FTChinese/superyard/pkg/fetch"
	"log"
	"net/http"
)

func (c Client) ConfirmOrder(orderID string) (*http.Response, error) {
	url := fmt.Sprintf("%s/orders/%s/verify-payment", c.baseURL, orderID)

	log.Printf("Query order payment result at %s", url)

	resp, errs := fetch.New().Patch(url).SetBearerAuth(c.key).End()
	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
