package torob

var CurrentRuntimeInfo = RuntimeInfo{
	MaxRunningWorkers: 1,
	WorkerPool: make(chan int, 1),
}
