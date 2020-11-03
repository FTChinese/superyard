package subsapi

import (
	"fmt"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/subs"
	"log"
	"net/http"
)

func (c Client) QueryOrder(order subs.Order) (*http.Response, error) {
	url := fmt.Sprintf("%s/orders/%s", c.baseURL, order.ID)

	log.Printf("Query order payment result at %s", url)

	resp, errs := fetch.New().Patch(url).SetBearerAuth(c.key).End()
	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
