package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/labstack/gommon/log"
)

func (env Env) SaveMemberSnapshot(s reader.MemberSnapshot) error {
	_, err := env.db.NamedExec(
		reader.InsertMemberSnapshot,
		s)
	if err != nil {
		return err
	}

	return nil
}

func (env Env) countMemberSnapshot(ids reader.IDs) (int64, error) {
	var count int64
	err := env.db.Get(&count, reader.StmtCountMemberSnapshot, ids.BuildFindInSet())
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listMemberSnapshot(ids reader.IDs, p gorest.Pagination) ([]reader.MemberSnapshot, error) {
	snapshots := make([]reader.MemberSnapshot, 0)

	err := env.db.Select(
		&snapshots,
		reader.StmtMemberSnapshots,
		ids.BuildFindInSet(),
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (env Env) ListMemberSnapshots(ids reader.IDs, p gorest.Pagination) (reader.MemberRevisions, error) {
	countCh := make(chan int64)
	listCh := make(chan reader.MemberRevisions)

	go func() {
		defer close(countCh)
		n, err := env.countMemberSnapshot(ids)
		if err != nil {
			log.Error(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)
		list, err := env.listMemberSnapshot(ids, p)
		listCh <- reader.MemberRevisions{
			Total:      0,
			Pagination: gorest.Pagination{},
			Data:       list,
			Err:        err,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return reader.MemberRevisions{}, listResult.Err
	}

	listResult.Total = count
	listResult.Pagination = p

	return listResult, nil
}
