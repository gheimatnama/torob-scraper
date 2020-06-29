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
	err := getJson(fmt.Sprintf("https://gheimatnama.ir/products/listNames?page=%d", page), &pagination, false)
	if err != nil {
		log.Error(err)
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
	_, lastPage := ParseSearchQueries(1)
	for page := 1; page <= lastPage; page++ {
		queries, _ := ParseSearchQueries(page)
		queries = filterExistedSearchQueries(queries)
		log.Info("Filtered search queries - only ", len(queries), " are valid")
		SearchAndPersist(pluckTitleFromSearchQueries(queries))
		PersistSearchQueries(queries)
	}
	log.Info("All done!")
}