package torob


func PersistProduct(product *Product) {
	CurrentRuntimeInfo.DB.Save(product)
	sources := product.ProductsInfo.Result
	tx := CurrentRuntimeInfo.DB.Begin()
	for i := range sources {
		if !SourceExists(product, &sources[i]){
			CurrentRuntimeInfo.DB.Create(&sources[i])
		}
	}
	tx.Commit()
}


func PersistSearchQuery(query *SearchQuery) {
	CurrentRuntimeInfo.DB.Save(query)
}

func PersistSearchQueries(queries []SearchQuery) {
	for _, query := range queries {
		PersistSearchQuery(&query)
	}
}

func ProductExistsByID(ID string) bool {
	count := 0
	CurrentRuntimeInfo.DB.Model(&Product{}).Where("random_key = ?", ID).Count(&count)
	return count != 0
}


func SearchQueryExistsByProductID(ID int) bool {
	count := 0
	CurrentRuntimeInfo.DB.Model(&SearchQuery{}).Where("product_id = ?", ID).Count(&count)
	return count != 0
}


func SourceExists(product *Product, source *ProductSource) bool {
	count := 0
	CurrentRuntimeInfo.DB.Model(&ProductSource{}).Where("product_id = ? AND shop_id = ?", product.RandomKey, source.ShopId).Count(&count)
	return count != 0
}


func ListFailedSources() []ProductSource {
	var failedSources []ProductSource
	CurrentRuntimeInfo.DB.Where("direct_page_url = ?", "").Find(&failedSources)
	return failedSources
}

func UpdateProductSource(productSource *ProductSource) {
	CurrentRuntimeInfo.DB.Save(productSource)
}