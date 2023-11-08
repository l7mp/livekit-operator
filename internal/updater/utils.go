package updater

import (
	"fmt"
	"github.com/l7mp/livekit-operator/internal/store"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (u *Updater) upsertService(svc *corev1.Service, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert service", "resource", store.GetObjectKey(svc), "generation", gen)

	mgrClient := u.manager.GetClient()
	current := &corev1.Service{ObjectMeta: metav1.ObjectMeta{
		Name:      svc.GetName(),
		Namespace: svc.GetNamespace(),
	}}

	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		if err := mergeMetadata(current, svc); err != nil {
			return nil
		}

		// rewrite spec
		svc.Spec.DeepCopyInto(&current.Spec)

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert service %q: %w",
			store.GetObjectKey(svc), err)
	}

	u.log.V(1).Info("service upserted", "resource", store.GetObjectKey(svc), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current)

	return op, nil
}

func (u *Updater) upsertDeployment(dp *appv1.Deployment, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert deployment", "resource", store.GetObjectKey(dp), "generation", gen)

	mgrClient := u.manager.GetClient()
	current := &v1.Deployment{ObjectMeta: metav1.ObjectMeta{
		Name:      dp.GetName(),
		Namespace: dp.GetNamespace(),
	}}

	op, err := ctrlutil.CreateOrPatch(u.ctx, mgrClient, current, func() error {
		if err := mergeMetadata(current, dp); err != nil {
			return nil
		}

		current.Spec.Selector = dp.Spec.Selector
		if dp.Spec.Replicas != nil {
			current.Spec.Replicas = dp.Spec.Replicas
		}

		dp.Spec.Template.ObjectMeta.DeepCopyInto(&current.Spec.Template.ObjectMeta)

		currentSpec := &current.Spec.Template.Spec
		dpSpec := &dp.Spec.Template.Spec

		currentSpec.Containers = make([]corev1.Container, len(dpSpec.Containers))
		for i := range dpSpec.Containers {
			dpSpec.Containers[i].DeepCopyInto(&currentSpec.Containers[i])
		}

		// rest is optional
		if dpSpec.TerminationGracePeriodSeconds != nil {
			currentSpec.TerminationGracePeriodSeconds = dpSpec.TerminationGracePeriodSeconds
		}

		currentSpec.HostNetwork = dpSpec.HostNetwork

		// affinity
		if dpSpec.Affinity != nil {
			currentSpec.Affinity = dpSpec.Affinity
		}

		// tolerations
		if dpSpec.Tolerations != nil {
			currentSpec.Tolerations = dpSpec.Tolerations
		}

		// security context
		if dpSpec.SecurityContext != nil {
			currentSpec.SecurityContext = dpSpec.SecurityContext
		}

		currentSpec.ServiceAccountName = dpSpec.ServiceAccountName

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert deployment %q: %w",
			store.GetObjectKey(dp), err)
	}

	return op, nil
}

func mergeMetadata(dst, src client.Object) error {
	labs := labels.Merge(dst.GetLabels(), src.GetLabels())
	dst.SetLabels(labs)

	annotations := dst.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	for k, v := range src.GetAnnotations() {
		annotations[k] = v
	}
	dst.SetAnnotations(annotations)

	return addOwnerRef(dst, src)
}

func addOwnerRef(dst, src client.Object) error {
	ownerRefs := src.GetOwnerReferences()
	if len(ownerRefs) != 1 {
		return fmt.Errorf("addOwnerRef: expecting a singleton ownerRef in %q, found %d",
			store.GetObjectKey(src), len(ownerRefs))
	}
	ownerRef := src.GetOwnerReferences()[0]

	for i, ref := range dst.GetOwnerReferences() {
		if ref.Name == ownerRef.Name && ref.Kind == ownerRef.Kind {
			ownerRefs = dst.GetOwnerReferences()
			ownerRef.DeepCopyInto(&ownerRefs[i])
			dst.SetOwnerReferences(ownerRefs)

			return nil
		}
	}

	ownerRefs = dst.GetOwnerReferences()
	ownerRefs = append(ownerRefs, ownerRef)
	dst.SetOwnerReferences(ownerRefs)

	return nil
}
