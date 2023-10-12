package event

import (
	"fmt"
)

// Type  specifies the type of event sent to the operator
type Type int

const (
	TypeRender Type = iota + 1
	TypeUpsert
	TypeDelete
	TypeUpdate
	TypeUnknown
)

const (
	eventTypeRenderStr = "render"
	eventTypeUpsertStr = "upsert"
	eventTypeDeleteStr = "delete"
	eventTypeUpdateStr = "update"
)

// NewEventType parses an event type specification
func NewEventType(raw string) (Type, error) {
	switch raw {
	case eventTypeRenderStr:
		return TypeRender, nil
	case eventTypeUpsertStr:
		return TypeUpsert, nil
	case eventTypeDeleteStr:
		return TypeDelete, nil
	case eventTypeUpdateStr:
		return TypeUpdate, nil
	default:
		return TypeUnknown, fmt.Errorf("unknown event type: %q", raw)
	}
}

// String returns a string representation for an event
func (a Type) String() string {
	switch a {
	case TypeRender:
		return eventTypeRenderStr
	case TypeUpsert:
		return eventTypeUpsertStr
	case TypeDelete:
		return eventTypeDeleteStr
	case TypeUpdate:
		return eventTypeUpdateStr
	default:
		return "<unknown>"
	}
}

// Event defines an event sent to/from the operator
type Event interface {
	GetType() Type
	String() string
}
