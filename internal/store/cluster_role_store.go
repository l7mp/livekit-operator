package store

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/types"
)

var ClusterRoles = NewClusterRoleStore()

type ClusterRoleStore struct {
	Store
}

func NewClusterRoleStore() *ClusterRoleStore {
	return &ClusterRoleStore{
		Store: NewStore(),
	}
}

// GetAll returns all ClusterRole objects from the global storage
func (s *ClusterRoleStore) GetAll() []*rbacv1.ClusterRole {
	ret := make([]*rbacv1.ClusterRole, 0)

	objects := s.Objects()
	for i := range objects {
		r, ok := objects[i].(*rbacv1.ClusterRole)
		if !ok {
			// this is critical: throw up hands and die
			panic("access to an invalid object in the global ClusterRoleStore")
		}

		ret = append(ret, r)
	}

	return ret
}

// GetObject returns a named ClusterRole object from the global storage
func (s *ClusterRoleStore) GetObject(nsName types.NamespacedName) *rbacv1.ClusterRole {
	o := s.Get(nsName)
	if o == nil {
		return nil
	}

	r, ok := o.(*rbacv1.ClusterRole)
	if !ok {
		// this is critical: throw up hands and die
		panic("access to an invalid object in the global ClusterRoleStore")
	}

	return r
}
