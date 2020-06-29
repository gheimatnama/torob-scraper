package torob

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)


// Returns search queries and last available page
func ParseSearchQueries(page int) ([]SearchQuery, int) {
	log.Info("Parsing titles page ", page)
	pagination := Pagination{}
	var queries []SearchQuery
	err := getJson(fmt.Sprintf("https://gheimatnama.ir/products/listNames?page=%d", page), &pagination)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range pagination.Data {
		queries = append(queries, SearchQuery{
			ProductID: int(item.(map[string]interface{})["id"].(float64)),
			Title:     item.(map[string]interface{})["title"].(string),
			CreatedAt: time.Now(),
		})
	}
	log.Info("Found ", len(queries), " titles to search")
	return queries, pagination.LastPage
}


func pluckTitleFromSearchQueries(queries []SearchQuery) []string {
	var titles []string
	for _, query := range queries {
		titles = append(titles, query.Title)
	}
	return titles
}


func filterExistedSearchQueries(queries []SearchQuery) []SearchQuery {
	var newQueries []SearchQuery
	for _, query := range queries {
		if !SearchQueryExistsByProductID(query.ProductID) {
			newQueries = append(newQueries, query)
		}
	}
	return newQueries
}

func ParseQueriesAndSearch() {
	log.Info("Starting with ", CurrentRuntimeInfo.MaxRunningWorkers, " workers!")
	queries, lastPage := ParseSearchQueries(1)
	for page := 2; page <= lastPage; page++ {
		queries := filterExistedSearchQueries(queries)
		SearchAndPersist(pluckTitleFromSearchQueries(queries))
		PersistSearchQueries(queries)
		queries, _ = ParseSearchQueries(page)
	}
	log.Info("All done!")
}