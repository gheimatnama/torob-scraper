package torob

var CurrentRuntimeInfo = RuntimeInfo{
	MaxRunningWorkers: 1,
	SearchResultLimit: 3,
	WorkerPool: make(chan int, 1),
}
