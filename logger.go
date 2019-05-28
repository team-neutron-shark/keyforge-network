package kfnetwork

import (
	"fmt"
	"sync"
)

var loggerOnce sync.Once
var loggerSingleton *LogManager

// Log - struct representing
type Log struct {
	message string
}

// LogManager - Shockingly enough this class manages logs.
type LogManager struct {
	logQueue   chan string
	errorQueue chan string
	limit      int
}

// Logger - Function used to access the LogManager singleton pointer.
func Logger() *LogManager {
	loggerOnce.Do(func() {
		loggerSingleton = new(LogManager)
		loggerSingleton.limit = 1024
		loggerSingleton.logQueue = make(chan string, loggerSingleton.limit)
	})

	return loggerSingleton
}

// GetQueue - Returns a pointer to the LogManager log queue channel.
func (l *LogManager) GetQueue() *chan string {
	return &l.logQueue
}

// Log - Adds a log in the form of a string to the log queue.
func (l *LogManager) Log(message string) {
	logMessage := fmt.Sprintf("[   LOG   ] %s", message)
	l.logQueue <- logMessage
}

// GetLogs - Pull logs off of the channel queue and return them in an array.
// This is primarily used in conjunction with PrintLogs().
func (l *LogManager) GetLogs() []string {
	logs := []string{}

	for len(l.logQueue) > 0 {
		logs = append(logs, <-l.logQueue)
	}

	return logs
}

// PrintLogs - Prints the logs waiting on the queue.
func (l *LogManager) PrintLogs() {
	for _, log := range l.GetLogs() {
		fmt.Println(log)
	}
}

func (l *LogManager) Error(message string) {
	logMessage := fmt.Sprintf("[  ERROR  ] %s", message)
	l.logQueue <- logMessage
}

func (l *LogManager) Warn(message string) {
	logMessage := fmt.Sprintf("[ WARNING ] %s", message)
	l.logQueue <- logMessage
}
