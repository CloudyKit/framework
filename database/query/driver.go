package query

type Query struct {
}

/*
	result := db.Search(ProductScheme,
		query.New(
			query.Eq("",""),
			query.Neq("",""),
			query.In("",""),
			query.NotIn("",""),
		).Limit(1).
		Order("-Price")
	 )

	result.Fetch(&product)
	result.Fetch(&product2)

*/
