package fta

import (
	"fmt"
	"github.com/FTChinese/superyard/pkg/fetch"
)

func (c Client) RefreshTermsDoc(id string) error {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(baseTerms).
		AddPath(id).
		AddQueryBool("refresh", true).
		String()

	resp, errs := fetch.New().Get(url).End()
	if errs != nil {
		return errs[0]
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}
