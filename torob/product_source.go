package torob

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	//"fmt"
	//"time"
)

func getSourceDirectUrl(source ProductSource)  (string, error) {
	url := source.PageUrl
	url = url[:strings.Index(url, "&")]
	url = url + "&source=next_pwa&uid=&discover_method=direct_next_pwa&experiment="
	data, err := getText(url)
	if err != nil {
		return "", err
	}
	str := `window.location.href="`
	index := strings.Index(data, str)
	if index < 0 {
		return "", errors.New("Can't find url , content : " + data)
	}
	index += len(str)
	to := strings.Index(data[index:], `";</script>`) + index
	return data[index:to], nil
}

func GetDigikalaSourceID(source *ProductSource) string {
	url := source.DirectPageUrl
	index := strings.LastIndex(url, "/") + 1
	realUrl := Base64Decode(url[index:])
	lastPart := realUrl[strings.Index(realUrl, "dkp-") + 4:]
	ID := lastPart
	if strings.Contains(lastPart, "/") {
		ID = lastPart[:strings.Index(lastPart, "/") + 1]
	}
	return ID
}

func FillIDInSource(source *ProductSource) {
	providersMap := make(map[string]func(*ProductSource)string)
	providersMap["affstat.adro.co"] = GetDigikalaSourceID
	providerFunc, ok := providersMap[GetHostName(source.DirectPageUrl)]
	if ok {
		source.IDInSource = providerFunc(source)
	}
}

func CheckForProductLinkPotential(product *Product) bool {
	for index := range product.ProductsInfo.Result {
		if product.ProductsInfo.Result[index].ShopId == 10 || strings.Contains(product.ProductsInfo.Result[index].ShopName, "دیجیکالا") {
			return true
		}
	}
	return false
}

func FillProductSources(product *Product) {
	if product.ProductsInfo.Result == nil {
		return
	}
	var wg sync.WaitGroup
	if !CheckForProductLinkPotential(product) {
		log.Info("Parsing product", product.RandomKey ,"with no linkability skipped")
		return
	}
	log.Info("Found ", len(product.ProductsInfo.Result), " sources for ", product.RandomKey)
	for index := range product.ProductsInfo.Result {
		productSource := &product.ProductsInfo.Result[index]
		productSource.ProductID = product.RandomKey
		productSource.SourceID = GetQueryParam(productSource.PageUrl, "prk")[0]
		wg.Add(1)
		go func(source *ProductSource, wg *sync.WaitGroup) {
			defer wg.Done()
			DownloadProductSource(source)
		}(productSource, &wg)
	}
	wg.Wait()
	log.Info("Parsed all sources for ", product.RandomKey)
}


func DownloadProductSource(productSource *ProductSource) {
	directUrl, err := getSourceDirectUrl(*productSource)
	if err == nil {
		productSource.DirectPageUrl = directUrl
		FillIDInSource(productSource)
		log.Info("Parsed source  ", productSource.DirectPageUrl, " for ", productSource.ProductID)
	} else {
		log.Error("Error parsing source  ", productSource.PageUrl, " for ", productSource.ProductID, " ", err)
	}
}

func ReDownloadFailedSources() {
	sources := ListFailedSources()
	log.Info("Found ", len(sources), " failed sources")
	var wg sync.WaitGroup
	pool := make(chan int, 30)
	for _, source := range sources {
		wg.Add(1)
		pool <- 1
		go func(source ProductSource, w *sync.WaitGroup) {
			defer wg.Done()
			DownloadProductSource(&source)
			UpdateProductSource(&source)
			if source.DirectPageUrl != "" {
				log.Info("Source ", source.ID, " updated")
			}
			<- pool
		}(source, &wg)
	}
	wg.Wait()
}
