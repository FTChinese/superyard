package stst

import (
	gorest "github.com/FTChinese/go-rest"
	stats2 "github.com/FTChinese/superyard/internal/pkg/stats"
	"log"
)

func (env Env) countAliUnconfirmed() (int64, error) {
	var count int64
	err := env.dbs.Read.Get(&count, stats2.StmtCountAliUnconfirmed)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listAliUnconfirmed(p gorest.Pagination) ([]stats2.UnconfirmedOrder, error) {
	orders := make([]stats2.UnconfirmedOrder, 0)

	err := env.dbs.Read.Select(
		&orders,
		stats2.StmtAliUnconfirmed,
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (env Env) AliUnconfirmed(p gorest.Pagination) (stats2.AliWxFailedList, error) {
	countCh := make(chan int64)
	listCh := make(chan stats2.AliWxFailedList)

	go func() {
		defer close(countCh)
		n, err := env.countAliUnconfirmed()
		if err != nil {
			log.Print(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)
		list, err := env.listAliUnconfirmed(p)
		listCh <- stats2.AliWxFailedList{
			Total:      0,
			Pagination: gorest.Pagination{},
			Data:       list,
			Err:        err,
		}
	}()

	count, listRes := <-countCh, <-listCh

	if listRes.Err != nil {
		return stats2.AliWxFailedList{}, listRes.Err
	}

	listRes.Total = count
	listRes.Pagination = p

	return listRes, nil
}

func (env Env) countWxUnconfirmed() (int64, error) {
	var count int64
	err := env.dbs.Read.Get(&count, stats2.StmtCountWxUnconfirmed)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listWxUnconfirmed(p gorest.Pagination) ([]stats2.UnconfirmedOrder, error) {
	orders := make([]stats2.UnconfirmedOrder, 0)

	err := env.dbs.Read.Select(
		&orders,
		stats2.StmtWxUnconfirmed,
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (env Env) WxUnconfirmed(p gorest.Pagination) (stats2.AliWxFailedList, error) {
	countCh := make(chan int64)
	listCh := make(chan stats2.AliWxFailedList)

	go func() {
		defer close(countCh)
		n, err := env.countWxUnconfirmed()
		if err != nil {
			log.Print(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)
		list, err := env.listWxUnconfirmed(p)
		listCh <- stats2.AliWxFailedList{
			Total:      0,
			Pagination: gorest.Pagination{},
			Data:       list,
			Err:        err,
		}
	}()

	count, listRes := <-countCh, <-listCh

	if listRes.Err != nil {
		return stats2.AliWxFailedList{}, listRes.Err
	}

	listRes.Total = count
	listRes.Pagination = p

	return listRes, nil
}
