package store

import (
	"encoding/json"
	"fmt"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
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

//func GetConfigMapsNamespacedNameFromLiveKitMesh(mesh *lkstnv1a1.LiveKitMesh) *types.NamespacedName {
//	return &types.NamespacedName{
//		Namespace: *mesh.Spec.Components.LiveKit.Deployment.ConfigMap.Namespace,
//		Name:      *mesh.Spec.Components.LiveKit.Deployment.ConfigMap.Name,
//	}
//}

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

// DumpObject converts an object into a human-readable form for logging.
func DumpObject(o client.Object) string {
	// default dump
	output := fmt.Sprintf("%#v", o)

	// copy
	ro := o.DeepCopyObject()

	var tmp client.Object
	switch ro := ro.(type) {
	case *lkstnv1a1.LiveKitMesh:
		tmp = ro
	case *corev1.Service:
		tmp = ro
	case *corev1.ConfigMap:
		//tmp = stripCM(ro)
		tmp = ro
	default:
		// this is not fatal
		return output
	}

	// remove cruft
	//tmp = strip(tmp)

	if marshaledJson, err := json.Marshal(tmp); err == nil {
		output = string(marshaledJson)
	}
	return output
}

func FetchAllObjectsBasedOnLabelFromAllStores(lkMeshName string) []client.Object {
	var objects []client.Object

	objects = append(objects, Services.FetchObjectBasedOnLabel(lkMeshName)...)
	objects = append(objects, Deployments.FetchObjectBasedOnLabel(lkMeshName)...)
	//objects = append(objects, Deployments.FetchObjectBasedOnLabel(lkMeshName)...)
	return objects
}
