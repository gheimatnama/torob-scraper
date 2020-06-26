package main

import (
	"encoding/json"
	"flag"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"time"
	"torobSpider/rotator"
	"torobSpider/torob"
	//"os"
)


func GetWorkersCount() *int {
	return flag.Int("workers", 1, "Total workers to scrape pages and images")
}

func GetQueriesFile() *string {
	return flag.String("queries", "queries.json", "Json array file containing all queries")
}

func OnlyRepairDownloadedSources() *bool {
	return flag.Bool("repair", false, "Only repair already downloaded sources")
}

func ParseRuntimeInfo() {
	workersCount := GetWorkersCount()
	queriesFile := GetQueriesFile()
	onlyRepair := OnlyRepairDownloadedSources()
	flag.Parse()
	torob.CurrentRuntimeInfo.MaxRunningWorkers = *workersCount
	torob.CurrentRuntimeInfo.WorkerPool = make(chan int, *workersCount)
	torob.CurrentRuntimeInfo.QueriesFile = *queriesFile
	torob.CurrentRuntimeInfo.OnlyRepairDownloadedSources = *onlyRepair
}

func ParseQueries() []string {
	plan, _ := ioutil.ReadFile(torob.CurrentRuntimeInfo.QueriesFile)
	var data []string
	json.Unmarshal(plan, &data)
	return data
}

func GetRotator() *rotator.ProxyRotator {
	scyllaProvider := rotator.NewProviderInstance(rotator.ParseScyllaProxies, time.Second * 120)
	fileProvider := rotator.NewProviderInstance(rotator.ParseProxyFile, time.Second * 120)
	checker := rotator.NewCheckerInstance(rotator.RecaptchaChecker)
	rotator := rotator.NewInstance(checker, []*rotator.ProxyProvider{scyllaProvider, fileProvider})
	rotator.ParallelProxyConnection = false
	rotator.ProxyConnectionDelay = time.Second * 5
	rotator.CheckProxyInterval = time.Minute * 2
	rotator.CheckProxyBeforeConnection = false
	rotator.ProxyQueueRetryTimeout = time.Second * 1
	rotator.ProxyQueueTimeout = time.Minute * 5
	return rotator
}

func main() {
	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&torob.Product{})
	db.AutoMigrate(&torob.ProductSource{})
	torob.CurrentRuntimeInfo.DB = db
	torob.CurrentRuntimeInfo.ProxyRotator = GetRotator()
	torob.CurrentRuntimeInfo.ProxyRotator.Init(30)
	ParseRuntimeInfo()
	//torob.ReDownloadFailedSources()
	torob.SearchAndPersist(ParseQueries())
}
