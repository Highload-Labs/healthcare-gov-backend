package shared

import "math"

type Pagination struct {
	Limit      int
	PageNumber int
}

type Metadata struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalData   int64 `json:"total_data"`
	TotalPages  int   `json:"total_pages"`
}

func (p *Pagination) SettleValue() {
	if p.Limit == 0 {
		p.Limit = 10
	}

	if p.Limit > 50 {
		p.Limit = 50
	}

	if p.PageNumber == 0 {
		p.PageNumber = 1
	}
}

func (p *Pagination) CalculateOffset() int {
	return (p.PageNumber - 1) * p.Limit
}

func (p *Pagination) CalculateTotalPages(totalData int64) int {
	return int(math.Ceil(float64(totalData) / float64(p.Limit)))
}
