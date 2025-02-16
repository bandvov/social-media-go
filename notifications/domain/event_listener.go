package domain

type EventListener interface {
	Subscribe(channel string, handler func(string)) error
	Publish(channel string, message string) error
}
