package torob

import (
	log "github.com/sirupsen/logrus"
)

func ParseFromUrl(url string) *Product {
	product := Product{}
	err := getJson(url, &product)
	if err != nil {
		log.Error("Error while parsing product : " + url + " -- ", err)
	}
	product.MoreInfoUrl = CleanProductUrl(product.MoreInfoUrl)
	log.Info("Product " + product.RandomKey + " parsed!")
	return &product
}

func ParseProductAndPersist(url string) {
	url = CleanProductUrl(url)
	id := GetQueryParam(url, "prk")[0]
	if ProductExistsByID(id) {
		return
	}
	log.Info("Parsing product : " + id)
	product := ParseFromUrl(url)
	FillProductSources(product)
	PersistProduct(product)
	log.Info("Product " + id + " parsed and persisted")
}