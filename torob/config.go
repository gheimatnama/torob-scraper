package torob

var CurrentRuntimeInfo = RuntimeInfo{
	MaxRunningWorkers: 1,
	SearchResultLimit: 3,
	MaxParallelSearch: make(chan int, 1),
	MaxParallelProductPerSearch: make(chan int, 5),
	WorkerPool: make(chan int, 1),
}
