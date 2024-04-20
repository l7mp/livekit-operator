package store

import (
	"k8s.io/apimachinery/pkg/types"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var HTTPRoutes = NewHTTPRouteStore()

type HTTPRouteStore struct {
	Store
}

func NewHTTPRouteStore() *HTTPRouteStore {
	return &HTTPRouteStore{
		Store: NewStore(),
	}
}

// GetAll returns all HTTPRoute objects from the global storage
func (s *HTTPRouteStore) GetAll() []*gwapiv1.HTTPRoute {
	ret := make([]*gwapiv1.HTTPRoute, 0)

	objects := s.Objects()
	for i := range objects {
		r, ok := objects[i].(*gwapiv1.HTTPRoute)
		if !ok {
			// this is critical: throw up hands and die
			panic("access to an invalid object in the global HTTPRouteStore")
		}

		ret = append(ret, r)
	}

	return ret
}

// GetObject returns a named HTTPRoute object from the global storage
func (s *HTTPRouteStore) GetObject(nsName types.NamespacedName) *gwapiv1.HTTPRoute {
	o := s.Get(nsName)
	if o == nil {
		return nil
	}

	r, ok := o.(*gwapiv1.HTTPRoute)
	if !ok {
		// this is critical: throw up hands and die
		panic("access to an invalid object in the global HTTPRouteStore")
	}

	return r
}
