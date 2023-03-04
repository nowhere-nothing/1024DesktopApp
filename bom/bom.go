package bom

type Event int
type EventType string
type EventHandler func(EventType, Event)

const (
	DOMContentLoaded EventType = "DOMContentLoaded"
)

type BOM struct {
	eventHandler map[EventType][]EventHandler
}

func NewBOM() *BOM {
	return &BOM{
		eventHandler: make(map[EventType][]EventHandler),
	}
}

func (b *BOM) AddEventListener(t EventType, f EventHandler) {
	b.eventHandler[t] = append(b.eventHandler[t], f)
}

func (b *BOM) EmitEvent(t EventType) {
	for _, f := range b.eventHandler[t] {
		f(t, 0)
	}
}
