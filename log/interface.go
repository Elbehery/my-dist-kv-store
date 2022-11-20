package log

type TransactionLogger interface {
	WritePut(Key, Value string)
	WriteDelete(Key string)
	Err() <-chan error
	ReadEvents() (<-chan Event, <-chan error)
	Run()
}
