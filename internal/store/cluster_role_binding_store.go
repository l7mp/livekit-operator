package store

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/types"
)

var ClusterRoleBindings = NewClusterRoleBindingStore()

type ClusterRoleBindingStore struct {
	Store
}

func NewClusterRoleBindingStore() *ClusterRoleBindingStore {
	return &ClusterRoleBindingStore{
		Store: NewStore(),
	}
}

// GetAll returns all ClusterRoleBinding objects from the global storage
func (s *ClusterRoleBindingStore) GetAll() []*rbacv1.ClusterRoleBinding {
	ret := make([]*rbacv1.ClusterRoleBinding, 0)

	objects := s.Objects()
	for i := range objects {
		r, ok := objects[i].(*rbacv1.ClusterRoleBinding)
		if !ok {
			// this is critical: throw up hands and die
			panic("access to an invalid object in the global ClusterRoleBindingStore")
		}

		ret = append(ret, r)
	}

	return ret
}

// GetObject returns a named ClusterRoleBinding object from the global storage
func (s *ClusterRoleBindingStore) GetObject(nsName types.NamespacedName) *rbacv1.ClusterRoleBinding {
	o := s.Get(nsName)
	if o == nil {
		return nil
	}

	r, ok := o.(*rbacv1.ClusterRoleBinding)
	if !ok {
		// this is critical: throw up hands and die
		panic("access to an invalid object in the global ClusterRoleBindingStore")
	}

	return r
}
