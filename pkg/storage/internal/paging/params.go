package paging

import "strings"

type Page struct {
	Number int
	Size   int
	Sort   string
	SortBy string
}

// ApplyDefaults sets default values for a Page object (in place).
func (p *Page) ApplyDefaults() {
	if p.Number <= 0 {
		p.Number = 1
	}

	if p.Size <= 0 {
		p.Size = 20
	}

	p.SortBy = strings.ToLower(p.SortBy)
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}

	if strings.ToLower(p.Sort) == "asc" {
		p.Sort = "ASC"
	} else {
		p.Sort = "DESC"
	}
}

func (p *Page) Next() {
	p.Number++
}
