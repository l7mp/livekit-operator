package event

// render event

type Render struct {
	Type Type
}

// NewEventRender returns an event with the corresponding resource
func NewEventRender() *Render {
	return &Render{
		Type: TypeRender,
	}
}

func (e *Render) GetType() Type {
	return e.Type
}

func (e *Render) String() string {
	return e.Type.String()
}
