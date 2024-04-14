package store

import (
	"k8s.io/apimachinery/pkg/types"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var Gateways = NewGatewayStore()

type GatewayStore struct {
	Store
}

func NewGatewayStore() *GatewayStore {
	return &GatewayStore{
		Store: NewStore(),
	}
}

// GetAll returns all Gateway objects from the global storage
func (s *GatewayStore) GetAll() []*gwapiv1.Gateway {
	ret := make([]*gwapiv1.Gateway, 0)

	objects := s.Objects()
	for i := range objects {
		r, ok := objects[i].(*gwapiv1.Gateway)
		if !ok {
			// this is critical: throw up hands and die
			panic("access to an invalid object in the global GatewayStore")
		}

		ret = append(ret, r)
	}

	return ret
}

// GetObject returns a named Gateway object from the global storage
func (s *GatewayStore) GetObject(nsName types.NamespacedName) *gwapiv1.Gateway {
	o := s.Get(nsName)
	if o == nil {
		return nil
	}

	r, ok := o.(*gwapiv1.Gateway)
	if !ok {
		// this is critical: throw up hands and die
		panic("access to an invalid object in the global GatewayStore")
	}

	return r
}
