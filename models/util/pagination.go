package util

// NewPagination creates a new Pagination instance.
// p is the page number, r is the rows to retrieve.
func NewPagination(p, limit int64) Pagination {
	if p < 1 {
		p = 1
	}

	return Pagination{
		Page:  p,
		Limit: limit,
	}
}

// Pagination is used to calculate limit and offset parameter used int sql statement.
type Pagination struct {
	Page  int64 `query:"page"`     // Which page is requesting data.
	Limit int64 `query:"per_page"` // How many items per page.
}

func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.Limit < 1 {
		p.Limit = 20
	}
}

// Offset calculate the offset for SQL.
func (p Pagination) Offset() int64 {
	return (p.Page - 1) * p.Limit
}
