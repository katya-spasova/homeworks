package orderer_log_drainer

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func OrderedLogDrainer(logs chan (chan string)) (result chan string) {
	result = make(chan string, 100)
	m := make(map[int]chan string)
	go func() {
		var i = 0
		for {
			log, logsOpened := <-logs
			if !logsOpened {
				break
			}
			i += 1
			m[i] = make(chan string, 100)
			go func(num int) {
				for {
					logItem, logOpened := <-log
					if !logOpened {
						close(m[num])
						break
					}
					lineContents := []string{strconv.Itoa(num), logItem}
					var modifiedItem string = strings.Join(lineContents, "\t")
					m[num] <- modifiedItem
				}
			}(i)
		}
	}()

	go func() {
		num := 1
		maxTries := 3
		for {
			ch, ok := m[num]
			if ok {
				for {
					logItem, logOpened := <-ch
					if logOpened {
						result <- logItem
					} else {
						num += 1
						break
					}
				}
			} else {
				if maxTries > 0 {
					time.Sleep(10 * time.Millisecond)
					maxTries -= 1
				} else {
					break
				}
			}
		}
		close(result)
	}()
	return result
}

func ExampleWithTwoLogs() {
	logs := make(chan (chan string))
	orderedLog := OrderedLogDrainer(logs)

	log1 := make(chan string)
	logs <- log1
	log2 := make(chan string)
	logs <- log2
	close(logs)

	log1 <- "aaa"
	log2 <- "bbb"
	log1 <- "ccc"
	log2 <- "ddd"
	close(log1)
	close(log2)

	for logEntry := range orderedLog {
		fmt.Println(logEntry)
	}
	// Output:
	// 1	aaa
	// 1	ccc
	// 2	bbb
	// 2	ddd
}

func ExampleFromTaskDesctiption() {
	logs := make(chan (chan string))
	orderedLog := OrderedLogDrainer(logs)

	first := make(chan string)
	logs <- first
	second := make(chan string)
	logs <- second

	first <- "test message 1 in first"
	second <- "test message 1 in second"
	second <- "test message 2 in second"
	first <- "test message 2 in first"
	first <- "test message 3 in first"
	// Print the first message now just because we can
	fmt.Println(<-orderedLog)

	third := make(chan string)
	logs <- third

	third <- "test message 1 in third"
	first <- "test message 4 in first"
	close(first)
	second <- "test message 3 in second"
	close(third)
	close(logs)

	second <- "test message 4 in second"
	close(second)

	// Print all the rest of the messages
	for logEntry := range orderedLog {
		fmt.Println(logEntry)
	}
	// Output:
	// 1	test message 1 in first
	// 1	test message 2 in first
	// 1	test message 3 in first
	// 1	test message 4 in first
	// 2	test message 1 in second
	// 2	test message 2 in second
	// 2	test message 3 in second
	// 2	test message 4 in second
	// 3	test message 1 in third
}

func main() {
	ExampleWithTwoLogs()
	ExampleFromTaskDesctiption()
}
