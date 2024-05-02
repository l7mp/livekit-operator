package store

import (
	"k8s.io/apimachinery/pkg/types"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

var TCPRoutes = NewTCPRouteStore()

type TCPRouteStore struct {
	Store
}

func NewTCPRouteStore() *TCPRouteStore {
	return &TCPRouteStore{
		Store: NewStore(),
	}
}

// GetAll returns all TCPRoute objects from the global storage
func (s *TCPRouteStore) GetAll() []*gwapiv1a2.TCPRoute {
	ret := make([]*gwapiv1a2.TCPRoute, 0)

	objects := s.Objects()
	for i := range objects {
		r, ok := objects[i].(*gwapiv1a2.TCPRoute)
		if !ok {
			// this is critical: throw up hands and die
			panic("access to an invalid object in the global TCPRouteStore")
		}

		ret = append(ret, r)
	}

	return ret
}

// GetObject returns a named TCPRoute object from the global storage
func (s *TCPRouteStore) GetObject(nsName types.NamespacedName) *gwapiv1a2.TCPRoute {
	o := s.Get(nsName)
	if o == nil {
		return nil
	}

	r, ok := o.(*gwapiv1a2.TCPRoute)
	if !ok {
		// this is critical: throw up hands and die
		panic("access to an invalid object in the global TCPRouteStore")
	}

	return r
}
