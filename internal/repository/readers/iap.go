package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/apple"
)

func (env Env) ListIAP(p gorest.Pagination) ([]apple.Subscription, error) {
	s := make([]apple.Subscription, 0)

	err := env.db.Select(&s, apple.StmtListIAPSubs, p.Limit, p.Offset())
	if err != nil {
		return nil, err
	}

	return s, nil
}
