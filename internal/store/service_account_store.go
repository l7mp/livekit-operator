package store

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

var ServiceAccounts = NewServiceAccountStore()

type ServiceAccountStore struct {
	Store
}

func NewServiceAccountStore() *ServiceAccountStore {
	return &ServiceAccountStore{
		Store: NewStore(),
	}
}

// GetAll returns all ServiceAccount objects from the global storage
func (s *ServiceAccountStore) GetAll() []*corev1.ServiceAccount {
	ret := make([]*corev1.ServiceAccount, 0)

	objects := s.Objects()
	for i := range objects {
		r, ok := objects[i].(*corev1.ServiceAccount)
		if !ok {
			// this is critical: throw up hands and die
			panic("access to an invalid object in the global ServiceAccountStore")
		}

		ret = append(ret, r)
	}

	return ret
}

// GetObject returns a named ServiceAccount object from the global storage
func (s *ServiceAccountStore) GetObject(nsName types.NamespacedName) *corev1.ServiceAccount {
	o := s.Get(nsName)
	if o == nil {
		return nil
	}

	r, ok := o.(*corev1.ServiceAccount)
	if !ok {
		// this is critical: throw up hands and die
		panic("access to an invalid object in the global ServiceAccountStore")
	}

	return r
}
