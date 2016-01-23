package request_buffer

import (
	"fmt"
	"testing"
	"time"
)

// duplicated tests

//func TestNewRequesterAndStop(t *testing.T) {
//	fmt.Println("In TestNewRequesterAndStop")
//	var requester Requester = NewRequester(10, 10)
//	defer requester.Stop()
//	if requester == nil {
//		t.Errorf("the returned requester is nil")
//	}
//}
//
//func TestAddRequestRunsARequest(t *testing.T) {
//	fmt.Println("In TestAddRequestRunsARequest")
//	var requester = NewRequester(10, 10)
//	defer requester.Stop()
//	var ran = make(chan struct{})
//	fr := &fakeRequest{
//		id:        "fakeId",
//		cacheable: true,
//		run: func() (interface{}, error) {
//			close(ran)
//			return "result1", nil
//		},
//		setResult: okSetResult,
//	}
//	requester.AddRequest(fr)
//	<-ran
//}
//
//func TestNonCacheableRequests(t *testing.T) {
//	fmt.Println("In TestNonCacheableRequests")
//	var requester = NewRequester(10, 10)
//	defer requester.Stop()
//	var expected1, expected2 = "foo", "bar"
//	var setted = make(chan struct{})
//	fr := &fakeRequest{
//		id:        "fakeId",
//		cacheable: false,
//		run: func() (interface{}, error) {
//			return expected1, nil
//		},
//		setResult: nil,
//	}
//	requester.AddRequest(fr)
//	fr = &fakeRequest{
//		id:        "fakeId",
//		cacheable: false,
//		run: func() (interface{}, error) {
//			defer close(setted)
//			return expected2, nil
//		},
//		setResult: nil,
//	}
//	requester.AddRequest(fr)
//	<-setted
//}

func TestWithTimeout(t *testing.T) {
	fmt.Println("In TestWithTimeout")
	var requester = NewRequester(10, 10)
	defer requester.Stop()
	var expected1, expected2 = "foo", "bar"
	fr := &fakeRequest{
		id:        "fakeId",
		cacheable: false,
		run: func() (interface{}, error) {
			return expected1, nil
		},
		setResult: nil,
	}
	requester.AddRequest(fr)
	time.Sleep(1000 * time.Millisecond)
	fr = &fakeRequest{
		id:        "fakeId",
		cacheable: false,
		run: func() (interface{}, error) {
			return expected2, nil
		},
		setResult: nil,
	}
	requester.AddRequest(fr)
}

func TestCasheableWithTimeout(t *testing.T) {
	fmt.Println("In TestCasheableWithTimeout")
	var requester = NewRequester(10, 10)
	defer requester.Stop()
	var expected1, expected2 = "foo", "bar"
	fr := &fakeRequest{
		id:        "fakeId",
		cacheable: true,
		run: func() (interface{}, error) {
			return expected1, nil
		},
		setResult: nil,
	}
	requester.AddRequest(fr)
	time.Sleep(1000 * time.Millisecond)
	fr = &fakeRequest{
		id:        "fakeId",
		cacheable: true,
		run: func() (interface{}, error) {
			return expected2, nil
		},
		setResult: func(interface{}, error) {
			fmt.Println("Set result in fakeId3")
		},
	}
	requester.AddRequest(fr)
}

func TestCasheableAndWait(t *testing.T) {
	fmt.Println("In TestCasheableAndWait")
	var requester = NewRequester(10, 10)
	defer requester.Stop()
	var expected1, expected2 = "foo", "bar"
	fr := &fakeRequest{
		id:        "fakeId",
		cacheable: true,
		run: func() (interface{}, error) {
			return expected1, nil
		},
		setResult: nil,
	}
	requester.AddRequest(fr)
	fr = &fakeRequest{
		id:        "fakeId",
		cacheable: true,
		run: func() (interface{}, error) {
			return expected2, nil
		},
		setResult: func(interface{}, error) {
			fmt.Println("Set result in fakeId3")
		},
	}
	requester.AddRequest(fr)
	time.Sleep(1000 * time.Millisecond)
}

func TestCacheableRequests(t *testing.T) {
	fmt.Println("In TestCacheableRequests")
	var requester = NewRequester(10, 10)
	defer func() {
		fmt.Println("In defer")
		requester.Stop()
	}()
	var expected1, expected2 = "foo", "bar"
	var setted = make(chan struct{})
	fr := &fakeRequest{
		id:        "fakeId",
		cacheable: true,
		run: func() (interface{}, error) {
			return expected1, nil
		},
		setResult: func(interface{}, error) {
			fmt.Println("I'm calling set result")
			defer close(setted)
		},
	}
	requester.AddRequest(fr)
	fr = &fakeRequest{
		id:        "fakeId",
		cacheable: true,
		run: func() (interface{}, error) {
			return expected2, nil
		},
		setResult: func(interface{}, error) {
			fmt.Println("I'm calling set result")
			defer close(setted)
		},
	}
	requester.AddRequest(fr)
	<-setted
}

// bad test
//func TestFullCache(t *testing.T) {
//	fmt.Println("In TestFullCache")
//	var requester = NewRequester(1, 10)
//	var setted = make(chan struct{})
//	fr := &fakeRequest{
//		id:        "fakeId",
//		cacheable: true,
//		run: func() (interface{}, error) {
//			return "foo", nil
//		},
//		setResult: nil,
//	}
//	requester.AddRequest(fr)
//	fr = &fakeRequest{
//		id:        "fakeId2",
//		cacheable: true,
//		run: func() (interface{}, error) {
//			return "foo", nil
//		},
//		setResult: nil,
//	}
//	requester.AddRequest(fr)
//		fr = &fakeRequest{
//		id:        "fakeId3",
//		cacheable: true,
//		run: func() (interface{}, error) {
//			return "foo", nil
//		},
//		setResult: nil,
//	}
//	requester.AddRequest(fr)
//	fr = &fakeRequest{
//		id:        "fakeId4",
//		cacheable: true,
//		run: func() (interface{}, error) {
//			return "foo", nil
//		},
//		setResult: nil,
//	}
//	requester.AddRequest(fr)
//	fr = &fakeRequest{
//		id:        "fakeId",
//		cacheable: true,
//		run: func() (interface{}, error) {
//			return "foo", nil
//		},
//		setResult: nil,
//	}
//	requester.AddRequest(fr)
//	fr = &fakeRequest{
//		id:        "fakeId6",
//		cacheable: true,
//		run: func() (interface{}, error) {
//			defer close(setted)
//			return "bar", nil
//		},
//		setResult: nil,
//	}
//	requester.AddRequest(fr)
//	<-setted
//}

//func TestMoreCacheableRequests(t *testing.T) {
//	fmt.Println("In TestMoreCacheableRequests")
//	var requester = NewRequester(10, 10)
//	beginning := time.Now()
//	defer func () {
//		fmt.Println("In defer")
//		duration := time.Since(beginning)
//		if duration.Seconds() > 2 {
//			t.Error("Execution takes too much time")
//		}
//		requester.Stop()
//	}()
//	for i := 0; i < 1000; i++ {
//		fr := &fakeRequest{
//			id:        "fakeId",
//			cacheable: true,
//			run: func() (interface{}, error) {
//				time.Sleep(1 * time.Second)
//				return "foo", nil
//			},
//			setResult: func(interface{}, error) {
//				fmt.Println("I'm calling set result")
//			},
//		}
//		requester.AddRequest(fr)
//	}
//}

func TestMyThrottleSize(t *testing.T) {
	fmt.Println("In TestThrottleSize")
	var requester = NewRequester(10, 2)
	beginning := time.Now()
	defer func() {
		fmt.Println("In defer")
		duration := time.Since(beginning)
		fmt.Println("duration is ", duration.Seconds())
		requester.Stop()
	}()
	maxI := 100
	for i := 0; i < maxI; i++ {
		fr := &fakeRequest{
			id:        "fakeId",
			cacheable: true,
			run: func() (interface{}, error) {
				fmt.Println("Execiting fakeId")
				time.Sleep(1 * time.Second)
				return "foo", nil
			},
			setResult: func(interface{}, error) {
				fmt.Println("Set result for fakeId")
			},
		}
		requester.AddRequest(fr)
		fr = &fakeRequest{
			id:        "fakeId2",
			cacheable: true,
			run: func() (interface{}, error) {
				fmt.Println("Execute fakeId2")
				time.Sleep(1 * time.Second)
				return "foo", nil
			},
			setResult: func(interface{}, error) {
				fmt.Println("Set result in fakeId2")
			},
		}
		requester.AddRequest(fr)
		fr = &fakeRequest{
			id:        "fakeId3",
			cacheable: true,
			run: func() (interface{}, error) {
				fmt.Println("In fake3")
				time.Sleep(1 * time.Second)
				fmt.Println("In fake3 2")
				fmt.Println("I'm here")
				return "foo", nil
			},
			setResult: func(interface{}, error) {
				fmt.Println("Set result in fakeId3")
			},
		}
		requester.AddRequest(fr)
	}
}

func TestLongRequestAndStop(t *testing.T) {
	fmt.Println("In TestLongRequestAndStop")
	var requester = NewRequester(10, 2)
	fr := &fakeRequest{
		id:        "fakeId",
		cacheable: false,
		run: func() (interface{}, error) {
			defer func() {
				fmt.Println("End of run")
			}()
			fmt.Println("Executing fakeId")
			time.Sleep(2 * time.Second)
			return "foo", nil
		},
		setResult: func(interface{}, error) {
			fmt.Println("Set result for fakeId")
		},
	}
	requester.AddRequest(fr)

	fr = &fakeRequest{
		id:        "fakeId2",
		cacheable: false,
		run: func() (interface{}, error) {
			defer func() {
				fmt.Println("End of run fakeId2")
			}()
			fmt.Println("Executing fakeId2")
			time.Sleep(2 * time.Second)
			return "foo", nil
		},
		setResult: func(interface{}, error) {
			fmt.Println("Set result for fakeId")
		},
	}
	requester.AddRequest(fr)
	time.Sleep(1 * time.Millisecond)
	fmt.Println("before stop")
	requester.Stop()
	fmt.Println("Stop called")
}

type fakeRequest struct {
	id         string
	cacheable  bool
	alreadyRan bool
	run        func() (interface{}, error)
	setResult  func(interface{}, error)
}

func (fr *fakeRequest) ID() string {
	return fr.id
}

func (fr *fakeRequest) Run() (interface{}, error) {
	if fr.alreadyRan {
		panic("Run after Run or SetResult")
	}
	fr.alreadyRan = true
	return fr.run()
}

func (fr *fakeRequest) Cacheable() bool {
	return fr.cacheable
}

func (fr *fakeRequest) SetResult(result interface{}, err error) {
	if fr.alreadyRan {
		panic("SetResult after Run or SetResult")
	}
	fr.alreadyRan = true
	fr.setResult(result, err)
}

//duplicated func

//func okSetResult(result interface{}, err error) {
//}
