package store

import (
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetObjectKey(object client.Object) string {

	return types.NamespacedName{
		Namespace: object.GetNamespace(),
		Name:      object.GetName(),
	}.String()
}

func GetNamespacedName(object client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: object.GetNamespace(),
		Name:      object.GetName(),
	}
}

func GetConfigMapsNamespacedNameFromLiveKitMesh(mesh *lkstnv1a1.LiveKitMesh) *types.NamespacedName {
	return &types.NamespacedName{
		Namespace: mesh.GetNamespace(),
		Name:      *mesh.Spec.Components.LiveKit.Deployment.ConfigMap,
	}
}

// Two resources are different if:
// (1) They have different namespaces or names.
// (2) They have the same namespace and name (resources are the same resource) but their specs are different.
// If their specs are different, their Generations are different too. So we only test their Generations.
// note: annotations are not part of the spec, so their update doesn't affect the Generation.
func compareObjects(o1, o2 client.Object) bool {
	return o1.GetNamespace() == o2.GetNamespace() &&
		o1.GetName() == o2.GetName() &&
		o1.GetGeneration() == o2.GetGeneration()
}
