package store

import (
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
)

var LiveKitMeshes = NewLivekitMeshStore()

type LiveKitMeshStore struct {
	Store
}

func NewLivekitMeshStore() *LiveKitMeshStore {
	return &LiveKitMeshStore{
		Store: NewStore(),
	}
}

func (l *LiveKitMeshStore) IsConfigMapReadyForMesh(mesh *lkstnv1a1.LiveKitMesh) bool {
	storedConfigMaps := ConfigMaps.GetAll()
	for _, configMap := range storedConfigMaps {
		cm := configMap
		if GetNamespacedName(cm) == *GetConfigMapsNamespacedNameFromLiveKitMesh(mesh) {
			return true
		}
	}
	return false
}

func (l *LiveKitMeshStore) GetAll() []*lkstnv1a1.LiveKitMesh {
	ret := make([]*lkstnv1a1.LiveKitMesh, 0)

	objects := l.Objects()
	for i := range objects {
		r, ok := objects[i].(*lkstnv1a1.LiveKitMesh)
		if !ok {
			// this is critical: throw up hands and die
			panic("access to an invalid object in the global LiveKitMeshStore")
		}

		ret = append(ret, r)
	}

	return ret
}
