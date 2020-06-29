package torob

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"sync"
	"os"
	"bufio"
	"math/rand"
	"time"
)

func getSearchUrl(query string) string {
	return "https://api.torob.com/v4/base-product/search/?q="+
		url.QueryEscape(query) +
		"&sort=popularity&page=0&size=24&source=next_pwa"
}

func SearchByQuery(query string) []SearchItem {
	searchResult := SearchResult{}
	err := getJson(getSearchUrl(query), &searchResult, true)
	if err != nil {
		log.Error("Error while searching query : " + query + " -- ", err)
	}
	results := searchResult.Results
	if len(results) > CurrentRuntimeInfo.SearchResultLimit {
		results = searchResult.Results[0:CurrentRuntimeInfo.SearchResultLimit]
	}
	return results
}


func ParseAndPersistSearchProducts(items []SearchItem) {
	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		CurrentRuntimeInfo.MaxParallelProductPerSearch <- 1
		go func(url string, wg *sync.WaitGroup) {
			defer wg.Done()
			ParseProductAndPersist(url)
			<- CurrentRuntimeInfo.MaxParallelProductPerSearch
		}(item.MoreInfoUrl, &wg)
	}
	wg.Wait()
}

func RequestFaker() {
	urls := make([]string, 0)
	file, err := os.Open("fakeurls.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        urls = append(urls, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

		for {
			getFakeText(urls[rand.Intn(len(urls))])
			println("SENT FAKE REQ")
			time.Sleep((time.Duration(rand.Intn(20))) * time.Millisecond)
		}
}

func SearchAndPersist(queries []string) {
	var wg sync.WaitGroup
	for _, query := range queries {
		wg.Add(1)
		CurrentRuntimeInfo.MaxParallelSearch <- 1
		go func(query string, wg *sync.WaitGroup) {
			defer wg.Done()
			log.Info("Searching " + query)
			items := SearchByQuery(query)
			log.Info("Found ", len(items), " results for " + query)
			ParseAndPersistSearchProducts(items)
			log.Info("Persisted all results for search query : " + query)
			<- CurrentRuntimeInfo.MaxParallelSearch
		}(query, &wg)
	}
	wg.Wait()
}
