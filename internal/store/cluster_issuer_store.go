package store

import (
	cert "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"k8s.io/apimachinery/pkg/types"
)

var Issuers = NewIssuerStore()

type IssuerStore struct {
	Store
}

func NewIssuerStore() *IssuerStore {
	return &IssuerStore{
		Store: NewStore(),
	}
}

// GetAll returns all Issuer objects from the global storage
func (s *IssuerStore) GetAll() []*cert.Issuer {
	ret := make([]*cert.Issuer, 0)

	objects := s.Objects()
	for i := range objects {
		r, ok := objects[i].(*cert.Issuer)
		if !ok {
			// this is critical: throw up hands and die
			panic("access to an invalid object in the global IssuerStore")
		}

		ret = append(ret, r)
	}

	return ret
}

// GetObject returns a named Issuer object from the global storage
func (s *IssuerStore) GetObject(nsName types.NamespacedName) *cert.Issuer {
	o := s.Get(nsName)
	if o == nil {
		return nil
	}

	r, ok := o.(*cert.Issuer)
	if !ok {
		// this is critical: throw up hands and die
		panic("access to an invalid object in the global IssuerStore")
	}

	return r
}
