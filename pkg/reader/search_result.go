package reader

type SearchResult struct {
	ID string `db:"id"`
}

type SearchBy int

const (
	SearchByEmail SearchBy = iota
	SearchByWxName
)
