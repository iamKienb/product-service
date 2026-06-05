package service

type Query map[string]any

type QueryBuilder struct {
	must   []Query
	filter []Query
	should []Query
	sorts  []map[string]any
	from   int
	size   int
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		must:   []Query{},
		filter: []Query{},
		should: []Query{},
		sorts:  []map[string]any{},
		size:   10,
	}
}

func (b *QueryBuilder) WithPagination(from, size int) *QueryBuilder {
	b.from = from
	if size > 0 {
		b.size = size
	}
	return b
}

func (b *QueryBuilder) MustTerm(field string, value any) *QueryBuilder {
	b.must = append(b.must, Query{"term": Query{field: value}})
	return b
}

func (b *QueryBuilder) MustTerms(field string, values []any) *QueryBuilder {
	b.must = append(b.must, Query{"terms": Query{field: values}})
	return b
}

func (b *QueryBuilder) FilterTerm(field string, value any) *QueryBuilder {
	b.filter = append(b.filter, Query{"term": Query{field: value}})
	return b
}

func (b *QueryBuilder) FilterTerms(field string, values []string) *QueryBuilder {
	interfaceValues := make([]any, len(values))
	for i, v := range values {
		interfaceValues[i] = v
	}
	b.filter = append(b.filter, Query{"terms": Query{field: interfaceValues}})
	return b
}

func (b *QueryBuilder) ShouldTerm(field string, value any) *QueryBuilder {
	b.should = append(b.should, Query{"term": Query{field: value}})
	return b
}

func (b *QueryBuilder) ShouldTerms(field string, values []any) *QueryBuilder {
	b.should = append(b.should, Query{"terms": Query{field: values}})
	return b
}

func (b *QueryBuilder) WithSort(field string, order string) *QueryBuilder {
	b.sorts = append(b.sorts, map[string]any{
		field: map[string]any{"order": order},
	})
	return b
}

func (b *QueryBuilder) Nested(path string, childBuilder *QueryBuilder, innerHitName string) *QueryBuilder {
	if childBuilder == nil {
		return b
	}

	childResult := childBuilder.Build()
	childQuery, exists := childResult["query"]
	if !exists {
		return b
	}

	nestedBody := Query{
		"path":  path,
		"query": childQuery,
	}

	if innerHitName != "" {
		nestedBody["inner_hits"] = Query{
			"name": innerHitName,
		}
	}

	b.filter = append(b.filter, Query{"nested": nestedBody})
	return b
}

func (b *QueryBuilder) Build() map[string]any {
	boolQuery := map[string]any{}

	if len(b.must) > 0 {
		boolQuery["must"] = b.must
	}
	if len(b.filter) > 0 {
		boolQuery["filter"] = b.filter
	}
	if len(b.should) > 0 {
		boolQuery["should"] = b.should
	}

	if len(boolQuery) == 0 {
		boolQuery["must"] = []Query{{"match_all": map[string]any{}}}
	}

	body := map[string]any{
		"from":             b.from,
		"size":             b.size,
		"track_total_hits": true,
		"query": map[string]any{
			"bool": boolQuery,
		},
	}

	if len(b.sorts) > 0 {
		body["sort"] = b.sorts
	}

	return body
}
