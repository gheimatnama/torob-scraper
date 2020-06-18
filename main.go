package main

import (
	"encoding/json"
	"flag"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"torobSpider/torob"
	//"os"
)

const MaxRunningWorkers = 10

func GetWorkersCount() *int {
	return flag.Int("workers", 10, "Total workers to scrape pages and images")
}

func GetQueriesFile() *string {
	return flag.String("queries", "queries.json", "Json array file containing all queries")
}

func ParseRuntimeInfo() {
	workersCount := GetWorkersCount()
	queriesFile := GetQueriesFile()
	flag.Parse()
	torob.CurrentRuntimeInfo.MaxRunningWorkers = *workersCount
	torob.CurrentRuntimeInfo.WorkerPool = make(chan int, *workersCount)
	torob.CurrentRuntimeInfo.QueriesFile = *queriesFile
}

func ParseQueries() []string {
	plan, _ := ioutil.ReadFile(torob.CurrentRuntimeInfo.QueriesFile)
	var data []string
	json.Unmarshal(plan, &data)
	return data
}

func main() {
	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	//os.Setenv("https_proxy", "http://194.5.206.111:8081/")
	db.AutoMigrate(&torob.Product{})
	db.AutoMigrate(&torob.ProductSource{})
	torob.CurrentRuntimeInfo.DB = db
	ParseRuntimeInfo()
	torob.SearchAndPersist(ParseQueries())
}
