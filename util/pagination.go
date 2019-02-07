package util

type Pagination struct {
	PageNumber int64 // Which page is requesting data.
	RowCount   int64 // How many items per page.
}

// NewPagination creates a new Pagination instance.
func NewPagination(p, r int64) Pagination  {
	if p < 1 {
		p = 1
	}

	return Pagination{
		PageNumber: p,
		RowCount:   r,
	}
}

// Offset calculate the offset for SQL.
func (p Pagination) Offset() int64 {
	return (p.PageNumber - 1) * p.RowCount
}
