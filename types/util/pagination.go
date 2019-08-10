package util

type Pagination struct {
	page  int64 // Which page is requesting data.
	Limit int64 // How many items per page.
}

// NewPagination creates a new Pagination instance.
// p is the page number, r is the rows to retrieve.
func NewPagination(p, limit int64) Pagination {
	if p < 1 {
		p = 1
	}

	return Pagination{
		page:  p,
		Limit: limit,
	}
}

// Offset calculate the offset for SQL.
func (p Pagination) Offset() int64 {
	return (p.page - 1) * p.Limit
}
