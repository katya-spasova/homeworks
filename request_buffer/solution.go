package request_buffer

import (
	"sync"
)

type Request interface {
	// return the id if the Request. If two Requests are the same, they have the same id
	ID() string

	// Blocks while it executes the Request
	// Returns the result from the execution or an error
	// The result and the error shoukd not be passed to SetResult
	// of the current request - They are stored internally before they are returned
	Run() (result interface{}, err error)

	// Returns a flag if the Request is cacheble
	// Has an unidentified behaviour id the called before `Run`
	Cacheable() bool

	// Sets the result of the Request
	// Should not be called for Request for which `Run` has been called.
	SetResult(result interface{}, err error)
}

type Requester interface {
	// Adds request and executes it if this us necessary at the first possible time
	AddRequest(request Request)

	// Stops the Requester. Waits all started Requests and call `SetResult` for all requests
	// that for which a same type Request has been executed
	// No new Request should be added at this time. No Requsts should be queued for calling
	// of `SetResult`
	Stop()
}

// Returns a new Requester, which cashes the responses of cacheSize Requests
// and executed no more than throttleSize Requests at the same time
func NewRequester(cacheSize int, throttleSize int) Requester {
	requester := &MyRequester{cacheSize: cacheSize, throttleSize: throttleSize}
	requester.init()
	go requester.start()

	return requester
}

// Implemenation of Requester interface
type MyRequester struct {
	cacheSize     int
	throttleSize  int
	running       bool
	mutex         sync.Mutex
	queue         []Request
	cache         map[string]ExecutionResult
	cachedIds     []string
	executionPool map[string]*Request
	cond          *sync.Cond
	finishCond    chan (struct{})
}

// Initialises all fields of MyRequester
func (requester *MyRequester) init() {
	requester.running = true
	requester.mutex = sync.Mutex{}
	requester.queue = make([]Request, 0)
	requester.cache = make(map[string]ExecutionResult, 0)
	requester.cachedIds = make([]string, 0)
	requester.executionPool = make(map[string]*Request)
	condMutex := sync.Mutex{}
	condMutex.Lock()
	requester.cond = sync.NewCond(&condMutex)
	requester.finishCond = make(chan struct{})
}

// Locks the requester
func (requester *MyRequester) Lock() {
	requester.mutex.Lock()
}

// Unlocks the requester
func (requester *MyRequester) Unlock() {
	requester.mutex.Unlock()
}

// Adds a Request for execution. It will be executed if necessary at the first possible time
func (requester *MyRequester) AddRequest(request Request) {
	requester.Lock()
	defer requester.Unlock()
	if requester.running {
		requester.queue = append(requester.queue, request)
		requester.cond.Signal()
	}
}

func (requester *MyRequester) hasNoRequests() bool {
	return len(requester.queue) == 0 && len(requester.executionPool) == 0
}

// Stops MyRequester. All pending requests will be executed
func (requester *MyRequester) Stop() {
	requester.Lock()
	if requester.running {
		requester.running = false
	}
	requester.cond.Signal()
	if !requester.hasNoRequests() {
		requester.Unlock()
		<-requester.finishCond
	} else {
		requester.Unlock()
	}
}

// Waits for Requests and executes them or takes the result from the cache
func (requester *MyRequester) start() {
	for {
		requester.Lock()
		hasNoRequests := requester.hasNoRequests()
		if !requester.running && hasNoRequests {
			requester.Unlock()
			close(requester.finishCond)
			break
		} else if len(requester.queue) == 0 {
			requester.Unlock()
			// solution's weak point - unlocking before wait leads to deadlocks sometimes
			requester.cond.Wait()
		} else {
			requester.Unlock()
		}

		requester.executeRequest()
	}
}

// finds  the first available request and executes it
func (requester *MyRequester) executeRequest() {
	requester.Lock()
	defer requester.Unlock()
	for i := 0; i < len(requester.queue); i++ {
		request := requester.queue[i]
		id := request.ID()
		// check if it is cached
		executionResult, ok := requester.cache[id]
		if ok {
			request.SetResult(executionResult.result, executionResult.err)
			requester.queue = append(requester.queue[:i], requester.queue[i+1:]...)
			break
		}

		// check if request of the same type is executed right now
		_, executedNow := requester.executionPool[id]
		if executedNow {
			continue
		}

		// request is not cached and is not executed right now
		// remove the request if the requester is stopped
		if !requester.running {
			requester.queue = append(requester.queue[:i], requester.queue[i+1:]...)
			break
		}

		// the requester is running - execute the request
		// add the request to the execution pool if possible
		if len(requester.executionPool) < requester.throttleSize {
			requester.executionPool[id] = &request
			requester.queue = append(requester.queue[:i], requester.queue[i+1:]...)
			go requester.doExecute(request)
			break
		}
	}
}

func (requester *MyRequester) doExecute(request Request) {
	result, err := request.Run()
	if request.Cacheable() {
		requester.addToCache(request.ID(), result, err)
	}
	requester.Lock()
	defer requester.Unlock()

	// remove the request from the execution pool
	delete(requester.executionPool, request.ID())
	if !requester.running && requester.hasNoRequests() {
		requester.cond.Signal()
	}
}

// Adds the result of Request's execution to the cache
func (requester *MyRequester) addToCache(id string, result interface{}, err error) {
	executionResult := ExecutionResult{result: result, err: err}
	requester.Lock()
	defer requester.Unlock()
	if len(requester.cachedIds) == requester.cacheSize {
		removeId := requester.cachedIds[0]
		requester.cachedIds = requester.cachedIds[1:]
		delete(requester.cache, removeId)
	}
	requester.cachedIds = append(requester.cachedIds, id)
	requester.cache[id] = executionResult
}

type ExecutionResult struct {
	result interface{}
	err    error
}
