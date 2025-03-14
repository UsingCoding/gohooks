package notifications

func NewDummy() Service {
	return &dummy{}
}

type dummy struct{}

func (d dummy) Send(_ string) error { return nil }
