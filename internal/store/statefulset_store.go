package store

import (
	appv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/types"
)

var StatefulSets = NewStatefulSetStore()

type StatefulSetStore struct {
	Store
}

func NewStatefulSetStore() *StatefulSetStore {
	return &StatefulSetStore{
		Store: NewStore(),
	}
}

// GetAll returns all Stateful objects from the global storage
func (s *StatefulSetStore) GetAll() []*appv1.StatefulSet {
	ret := make([]*appv1.StatefulSet, 0)

	objects := s.Objects()
	for i := range objects {
		r, ok := objects[i].(*appv1.StatefulSet)
		if !ok {
			// this is critical: throw up hands and die
			panic("access to an invalid object in the global StatefulSetStore")
		}

		ret = append(ret, r)
	}

	return ret
}

// GetObject returns a named Stateful object from the global storage
func (s *StatefulSetStore) GetObject(nsName types.NamespacedName) *appv1.StatefulSet {
	o := s.Get(nsName)
	if o == nil {
		return nil
	}

	r, ok := o.(*appv1.StatefulSet)
	if !ok {
		// this is critical: throw up hands and die
		panic("access to an invalid object in the global StatefulSetStore")
	}

	return r
}
