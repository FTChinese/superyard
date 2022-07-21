package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"log"
	"net/http"
)

func (c Client) VerifyOrder(orderID string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(rootPathOrders).
		AddPath(orderID).
		AddPath("verify-payment").
		String()

	log.Printf("Query order payment result at %s", url)

	resp, errs := fetch.New().Post(url).SetBearerAuth(c.key).End()
	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
