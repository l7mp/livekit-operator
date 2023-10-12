package event

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Upsert struct {
	Type   Type
	Object client.Object
	// Params map[string]string
}

// NewEventUpsert returns a new Upsert event
func NewEventUpsert(o client.Object) *Upsert {
	return &Upsert{Type: TypeUpsert, Object: o}
}

func (e *Upsert) GetType() Type {
	return e.Type
}

func (e *Upsert) String() string {
	return fmt.Sprintf("%s: %s/%s", e.Type.String(),
		e.Object.GetName(), e.Object.GetNamespace())
}
