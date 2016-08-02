package database

import "github.com/CloudyKit/framework/database/scheme"

var db IDB

type Promotion struct {
	Product     *Product
	ProductID   string
	Description string
	Expire      string
	Percentage  int
	Price       float64
}

type Product struct {
	ProductID  string
	Name       string
	Promotions []Promotion
}

var (
	ProductScheme   = scheme.New("products", "ProductID")
	PromotionScheme = scheme.New("product_promotions", "PromotionID")

	_ = scheme.Init(ProductScheme, func(def *scheme.Def) {
		def.Field("Name", scheme.String{255})
		def.RefsChildren("Promotions", PromotionScheme, "ProductID")
	})

	_ = scheme.Init(PromotionScheme, func(def *scheme.Def) {
		def.Field("ProductID", scheme.String{})
		def.Field("Description", scheme.String{})
		def.Field("Expire", scheme.DateTime{})

		def.RefsParent("Product", ProductScheme, "ProductID")
	})
)

func ExampleDB() {

	var product = struct {
		ProductID string
		Name      string
	}{}

	result := db.Search(ProductScheme, nil)
	for i := 0; i < result.NumOfRecords(); i++ {
		err := result.Fetch(&product)
		if err != nil {
			// unexpected err
			panic(err)
		}
	}

	result = db.Search(ProductScheme, nil)
	for result.FetchNext(&product) {

	}

	var products []Product
	result.FetchAll(&products)

}
