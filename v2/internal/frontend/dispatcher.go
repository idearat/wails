package frontend

type Dispatcher interface {
	ProcessMessage(message string) error
}
