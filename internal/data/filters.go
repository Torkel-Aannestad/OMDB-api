package data

import (
	"strings"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func (f *Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if safeValue == f.Sort {
			return strings.TrimSuffix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}
func (f *Filters) sortDirection() string {
	if strings.Contains(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.PermitterValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
