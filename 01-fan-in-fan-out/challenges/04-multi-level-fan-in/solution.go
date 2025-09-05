package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type LogType string

const (
	ERROR_LOG_TYPE = LogType("ERROR")
	WARN_LOG_TYPE  = LogType("WARN")
	INFO_LOG_TYPE  = LogType("INFO")
)

const GROUPS = 3
const SERVICE_PER_GROUP = 4
const PROCESSOR_PER_GROUP = 3
const MIN_LOGS_TO_PRODUCE = 5

var allTypes = []LogType{INFO_LOG_TYPE, WARN_LOG_TYPE, ERROR_LOG_TYPE}

type LogMessage struct {
	ServiceName string
	Message     string
}

type Classification struct {
	Type LogType
	Log  LogMessage
}

func createRandomLogMessage(serviceName string) LogMessage {
	logTypeIndex := rand.Intn(len(allTypes))
	logType := allTypes[logTypeIndex]

	return LogMessage{ServiceName: serviceName, Message: fmt.Sprintf("%s: some log message", string(logType))}
}

func runService(serviceName string, output chan<- LogMessage, wg *sync.WaitGroup) {
	numOfLogsToProduce := rand.Intn(5) + MIN_LOGS_TO_PRODUCE

	for i := 0; i < numOfLogsToProduce; i++ {
		time.Sleep(time.Duration(rand.Int63n(300) * int64(time.Millisecond)))
		output <- createRandomLogMessage(serviceName)
	}

	wg.Done()
}

func process(input <-chan LogMessage, output chan<- Classification, wg *sync.WaitGroup) {
	for logMessage := range input {
		messageSegments := strings.Split(logMessage.Message, ":")
		typeString := messageSegments[0]
		classification := Classification{}
		classification.Log = logMessage

		if typeString == string(ERROR_LOG_TYPE) {
			classification.Type = ERROR_LOG_TYPE
		} else if typeString == string(WARN_LOG_TYPE) {
			classification.Type = WARN_LOG_TYPE
		} else if typeString == string(INFO_LOG_TYPE) {
			classification.Type = INFO_LOG_TYPE
		} else {
			panic(fmt.Sprintf("Unknown log type: %s", typeString))
		}

		output <- classification
	}

	wg.Done()
}

func collectLocally(input <-chan Classification, output chan<- map[string]int, wg *sync.WaitGroup) {
	resultsMap := make(map[string]int)

	for classification := range input {
		if classification.Type == ERROR_LOG_TYPE {
			resultsMap[classification.Log.ServiceName]++
		} else {
			resultsMap[classification.Log.ServiceName] += 0
		}
	}

	output <- resultsMap

	wg.Done()
}

func collectGlobally(input <-chan map[string]int, output chan<- map[string]int) {
	resultsMap := make(map[string]int)

	for classifications := range input {
		for key, count := range classifications {
			resultsMap[key] += count
		}
	}

	output <- resultsMap

	close(output)
}

func main() {
	localWg := &sync.WaitGroup{}
	localCh := make(chan map[string]int)
	globalCh := make(chan map[string]int)

	go collectGlobally(localCh, globalCh)

	for i := 0; i < GROUPS; i++ {
		serviceWg := &sync.WaitGroup{}
		serviceCh := make(chan LogMessage)
		processorWg := &sync.WaitGroup{}
		processorCh := make(chan Classification)

		localWg.Add(1)
		go collectLocally(processorCh, localCh, localWg)

		for j := 0; j < PROCESSOR_PER_GROUP; j++ {
			processorWg.Add(1)
			go process(serviceCh, processorCh, processorWg)
		}

		for k := 0; k < SERVICE_PER_GROUP; k++ {
			serviceWg.Add(1)
			go runService(fmt.Sprintf("service-%d-%d", i, k), serviceCh, serviceWg)
		}

		go func() {
			processorWg.Wait()
			close(processorCh)
		}()

		go func() {
			serviceWg.Wait()
			close(serviceCh)
		}()
	}

	go func() {
		localWg.Wait()
		close(localCh)
	}()

	results := <-globalCh

	for service, errorsCount := range results {
		fmt.Printf("%s: %d\n", service, errorsCount)
	}

	fmt.Println("Done")
}
