package log

import (
	"bufio"
	"fmt"
	"os"
)

const (
	PutEvent EventType = iota
	DeleteEvent
)

type EventType byte

type Event struct {
	Sequence  uint
	EventType EventType
	Key       string
	Value     string
}

type FileTransactionLogger struct {
	file         *os.File
	events       chan<- Event
	errors       <-chan error
	lastSequence uint
}

func (l *FileTransactionLogger) WritePut(Key, Value string) {
	l.events <- Event{
		EventType: PutEvent,
		Key:       Key,
		Value:     Value,
	}
}

func (l *FileTransactionLogger) WriteDelete(Key string) {
	l.events <- Event{
		EventType: DeleteEvent,
		Key:       Key,
	}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *FileTransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for e := range events {
			l.lastSequence++
			_, err := fmt.Fprintf(l.file, "%d\t%d\t%s\t%s\n", l.lastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
			l.file.Sync()
		}
	}()
}

func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)

	outEvents := make(chan Event)
	outErrors := make(chan error, 1)

	go func() {
		var e Event

		for scanner.Scan() {
			line := scanner.Text()
			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s", &e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outErrors <- fmt.Errorf("input parse error: %w", err)
				return
			}
			if l.lastSequence >= e.Sequence {
				outErrors <- fmt.Errorf("transaction number is out of Sequence")
				return
			}
			l.lastSequence = e.Sequence
			outEvents <- e
		}

		if err := scanner.Err(); err != nil {
			outErrors <- fmt.Errorf("transaction log read failed: %w", err)
			return
		}
	}()

	return outEvents, outErrors
}

func NewFileTransactionLogger(fileName string) (TransactionLogger, error) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("could not open log file: %w", err)
	}

	return &FileTransactionLogger{
		file:         f,
		events: make(chan Event),
		errors: make(chan error),
	}, nil
}
