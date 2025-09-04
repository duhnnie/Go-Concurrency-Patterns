package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type LogMessageType string

const (
	ErrorLogMessageType = LogMessageType("ERROR")
	InfoLogMessageType  = LogMessageType("INFO")
	WarnLogMessageType  = LogMessageType("WARN")
)

type LogMessage struct {
	Service string
	Message string
}

type LogClassification struct {
	Type LogMessageType
	Log  LogMessage
}

var types = []LogMessageType{ErrorLogMessageType, InfoLogMessageType, WarnLogMessageType}

func runService(serviceName string, logChannel chan<- LogMessage, wg *sync.WaitGroup) {
	messagesToGenerate := rand.Intn(3) + 5

	for i := 0; i < messagesToGenerate; i++ {
		messageType := types[rand.Intn(len(types))]
		message := fmt.Sprintf("%s: some log message", string(messageType))
		time.Sleep(time.Duration(rand.Intn(300) * int(time.Millisecond)))
		logChannel <- LogMessage{serviceName, message}
	}

	wg.Done()
}

func process(input <-chan LogMessage, output chan<- LogClassification, wg *sync.WaitGroup) {
	for logMessage := range input {
		splittedMessage := strings.Split(logMessage.Message, ":")
		firstSegment := splittedMessage[0]
		var messageType LogMessageType

		switch firstSegment {
		case string(ErrorLogMessageType):
			messageType = ErrorLogMessageType
		case string(InfoLogMessageType):
			messageType = InfoLogMessageType
		case string(WarnLogMessageType):
			messageType = WarnLogMessageType
		default:
			panic(fmt.Sprintf("unknown %s\n", firstSegment))
		}

		output <- LogClassification{messageType, logMessage}
	}

	wg.Done()
}

func collect(input <-chan LogClassification, output chan<- map[string]int) {
	resultMap := make(map[string]int)

	for classification := range input {
		if classification.Type == ErrorLogMessageType {
			resultMap[classification.Log.Service] += 1
		} else {
			resultMap[classification.Log.Service] += 0
		}
	}

	output <- resultMap
	close(output)
}

func main() {
	logChannel := make(chan LogMessage)
	classificationChannel := make(chan LogClassification)
	result := make(chan map[string]int)
	serviceWg := &sync.WaitGroup{}
	processsorWg := &sync.WaitGroup{}

	for i := 0; i < 3; i++ {
		serviceWg.Add(1)
		go runService(fmt.Sprintf("service%d", i+1), logChannel, serviceWg)
	}

	go func() {
		serviceWg.Wait()
		close(logChannel)
	}()

	for i := 0; i < 3; i++ {
		processsorWg.Add(1)
		go process(logChannel, classificationChannel, processsorWg)
	}

	go func() {
		processsorWg.Wait()
		close(classificationChannel)
	}()

	go collect(classificationChannel, result)
	fmt.Println(<-result)
}
