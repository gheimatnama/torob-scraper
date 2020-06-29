package main

import (
	"flag"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
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

func GetRotator() *rotator.ProxyRotator {
	scyllaProvider := rotator.NewProviderInstance(rotator.ParseScyllaProxies, time.Minute * 2)
	fileProvider := rotator.NewProviderInstance(rotator.ParseProxyFile, time.Minute * 2)
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
	db, err := gorm.Open("mysql", "root:root@/torobspider?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		//log.Fatal(err)
		panic("failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&torob.Product{})
	db.AutoMigrate(&torob.ProductSource{})
	db.AutoMigrate(&torob.SearchQuery{})
	torob.CurrentRuntimeInfo.DB = db
	torob.CurrentRuntimeInfo.ProxyRotator = GetRotator()
	requiredAliveProxies := 50
	logrus.Info("Looking for ", requiredAliveProxies, " alive proxies before start")
	torob.CurrentRuntimeInfo.ProxyRotator.Init(requiredAliveProxies)
	ParseRuntimeInfo()
	if torob.CurrentRuntimeInfo.OnlyRepairDownloadedSources {
		torob.ReDownloadFailedSources()
	} else {
		torob.ParseQueriesAndSearch()
	}
}
