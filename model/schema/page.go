package schema

// Page is utilized whenever a collection is requested to provide
// metadata alongside requested data in a paginated form.
type Page[Entity any] struct {
	Contents      []Entity `json:"contents"`
	Offset        int      `json:"offset"`
	Total         int      `json:"total"`
	FilteredTotal int      `json:"filteredTotal"`
}

// TODO: Document in v1 json
