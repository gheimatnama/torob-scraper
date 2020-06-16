package torob

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"sync"
)

func getSearchUrl(query string) string {
	return "https://api.torob.com/v4/base-product/search/?q="+
		url.QueryEscape(query) +
		"&sort=popularity&page=0&size=24&source=next_pwa";
}

func SearchByQuery(query string) []SearchItem {
	searchResult := SearchResult{}
	err := getJson(getSearchUrl(query), &searchResult);
	if err != nil {
		log.Error("Error while searching query : " + query + " -- ", err)
	}
	results := searchResult.Results
	return results
}


func ParseAndPersistSearchProducts(items []SearchItem) {
	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		go func(url string, wg *sync.WaitGroup) {
			defer wg.Done()
			ParseProductAndPersist(url)
		}(item.MoreInfoUrl, &wg)
	}
	wg.Wait()
}

func SearchAndPersist(queries []string) {
	var wg sync.WaitGroup
	log.Info("Starting with ", CurrentRuntimeInfo.MaxRunningWorkers, " workers!")
	for _, query := range queries {
		wg.Add(1)
		go func(query string, wg *sync.WaitGroup) {
			defer wg.Done()
			log.Info("Searching " + query)
			items := SearchByQuery(query)
			log.Info("Found ", len(items), " results for " + query)
			ParseAndPersistSearchProducts(items)
			log.Info("Persisted all results for search query : " + query)
		}(query, &wg)
	}
	wg.Wait()
	log.Info("All done!")
}