package event

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	// "sigs.k8s.io/controller-runtime/pkg/client"
)

// Kind specifies the Kind of object under deletion
type Kind int

const (
	KindLiveKitMesh Kind = iota + 1
	KindUnknown
)

// String returns a string representation for an event
func (a Kind) String() string {
	switch a {
	case KindLiveKitMesh:
		return "LiveKitMesh"
	default:
		return "<unknown>"
	}
}

type Delete struct {
	Type Type
	Kind Kind
	Key  types.NamespacedName
}

// NewEventDelete returns a Delete event
func NewEventDelete(kind Kind, key types.NamespacedName) *Delete {
	return &Delete{Type: TypeDelete, Kind: kind, Key: key}
}

func (e *Delete) GetType() Type {
	return e.Type
}

func (e *Delete) String() string {
	return fmt.Sprintf("%s: %s of type %s", e.Type.String(), e.Key, e.Kind.String())
}
