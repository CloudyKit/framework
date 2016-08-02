package scheme_test

import (
	"github.com/CloudyKit/framework/database/scheme"
)

var (
	ProductScheme   = scheme.New("products", "ProductID")
	CategoryScheme  = scheme.New("categories", "CategoryID")
	PromotionScheme = scheme.New("promotions", "PromotionID")
	CustomerScheme  = scheme.New("customers", "CustomerID")

	_ = scheme.Init(ProductScheme, func(def *scheme.Def) {

		def.Field("Name", scheme.String{255})
		def.Field("Age", scheme.Int{})

		def.Refs("Categories", CategoryScheme, "CategoryID", "products_to_categories")
		def.RefsChildren("Promotions", PromotionScheme, "PromotionID")

	})

	_ = scheme.Init(CategoryScheme, func(def *scheme.Def) {

		def.Field("Name", scheme.String{255})
		def.RefsFrom("Products", ProductScheme, "ProductID")
	})
)

func ExampleScheme() {

	//var productCat []struct {
	//	ProductId  string
	//	Name       string
	//	Categories []struct{ Name string }
	//}
	//
	//numofrows, err := conn.Search(ProductScheme, query.Query{}, &productCat)
	//
	//if err != nil {
	//	 something went wrong
	//println("something went wrong: ", err.Error())
	//return
	//}

	//println("Found ", numofrows, "records")

	//for _, product := range productCat {
	//	println("PRODUCT:", product.ProductId, " => ", product.Name)
	//	for _, category := range product.Categories {
	//		println("\tCATEGORY:", category.Name)
	//	}
	//}

}
