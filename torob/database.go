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


func ProductExistsByID(ID string) bool {
	count := 0
	CurrentRuntimeInfo.DB.Model(&Product{}).Where("random_key = ?", ID).Count(&count)
	return count != 0
}


func SourceExists(product *Product, source *ProductSource) bool {
	count := 0
	CurrentRuntimeInfo.DB.Model(&ProductSource{}).Where("product_id = ? AND shop_id = ?", product.RandomKey, source.ShopId).Count(&count)
	return count != 0
}