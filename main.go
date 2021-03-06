package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
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

func GetMinimumRequiredAliveProxy() *int {
	return flag.Int("required-proxies", 10, "Minimum required proxies to start")
}

func ParseRuntimeInfo() {
	workersCount := GetWorkersCount()
	queriesFile := GetQueriesFile()
	onlyRepair := OnlyRepairDownloadedSources()
	requiredProxies := GetMinimumRequiredAliveProxy()
	flag.Parse()
	torob.CurrentRuntimeInfo.MaxRunningWorkers = *workersCount
	torob.CurrentRuntimeInfo.WorkerPool = make(chan int, *workersCount)
	torob.CurrentRuntimeInfo.QueriesFile = *queriesFile
	torob.CurrentRuntimeInfo.OnlyRepairDownloadedSources = *onlyRepair
	torob.CurrentRuntimeInfo.MinimumRequiredAliveProxy = *requiredProxies
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
	ParseRuntimeInfo()
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
	logrus.Info("Looking for ", torob.CurrentRuntimeInfo.MinimumRequiredAliveProxy, " alive proxies before start")
	torob.CurrentRuntimeInfo.ProxyRotator.Init(torob.CurrentRuntimeInfo.MinimumRequiredAliveProxy)
	if torob.CurrentRuntimeInfo.OnlyRepairDownloadedSources {
		torob.ReDownloadFailedSources()
	} else {
		torob.ParseQueriesAndSearch()
	}
}
