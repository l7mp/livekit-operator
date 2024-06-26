package updater

import (
	"fmt"
	cert "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	stnrgwv1 "github.com/l7mp/stunner-gateway-operator/api/v1"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

func (u *Updater) upsertConfigMap(cm *corev1.ConfigMap, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert configmap", "resource", store.GetObjectKey(cm), "generation", gen)

	mgrclient := u.manager.GetClient()
	current := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cm.GetName(),
			Namespace: cm.GetNamespace(),
		},
	}

	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrclient, current, func() error {
		err := mergeMetadata(current, cm)
		if err != nil {
			return err
		}

		err = addOwnerRef(current, cm)
		if err != nil {
			return err
		}

		current.Data = make(map[string]string)
		for k, v := range cm.Data {
			current.Data[k] = v
		}
		u.log.Info("current configmap", "cm", store.GetObjectKey(cm), "data", cm.Data)

		return nil
	})
	if err != nil {
		return ctrlutil.OperationResultNone, err
	}

	u.log.V(1).Info("config-map upserted", "resource", store.GetObjectKey(cm), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

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
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

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

		currentSpec := &current.Spec.Template.Spec
		dpSpec := &dp.Spec.Template.Spec

		dp.Spec.Template.ObjectMeta.DeepCopyInto(&current.Spec.Template.ObjectMeta)

		if current.Annotations[opdefault.RelatedConfigMapKey] != "" {
			cm := store.ConfigMaps.Get(types.NamespacedName{
				Namespace: current.Namespace,
				Name:      current.Annotations[opdefault.RelatedConfigMapKey],
			})
			if cm != nil {
				current.Spec.Template.Annotations[opdefault.DefaultConfigMapResourceVersionKey] = cm.GetResourceVersion()
			}
		}

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

	u.log.V(1).Info("deployment upserted", "resource", store.GetObjectKey(dp), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertSecret(s *corev1.Secret, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert cluster issuer secret", "resource", store.GetObjectKey(s), "generation", gen)

	mgrClient := u.manager.GetClient()
	current := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
		Name:      s.GetName(),
		Namespace: s.GetNamespace(),
	}}

	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		if err := mergeMetadata(current, s); err != nil {
			return nil
		}

		current.Data = map[string][]byte{}
		current.Data[opdefault.DefaultClusterIssuerSecretApiTokenKey] = s.Data[opdefault.DefaultClusterIssuerSecretApiTokenKey]
		current.Type = s.Type

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert cluster issuer secret %q: %w",
			store.GetObjectKey(s), err)
	}

	u.log.V(1).Info("cluster issuer secret upserted", "resource", store.GetObjectKey(s), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertIssuer(i *cert.Issuer, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert issuer", "resource", store.GetObjectKey(i), "generation", gen)

	mgrClient := u.manager.GetClient()
	current := &cert.Issuer{ObjectMeta: metav1.ObjectMeta{
		Name:      i.GetName(),
		Namespace: i.GetNamespace(),
	}}

	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		if err := mergeMetadata(current, i); err != nil {
			return nil
		}

		// rewrite spec
		i.Spec.DeepCopyInto(&current.Spec)

		u.log.Info("issuer", "i", i)

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert issuer %q: %w",
			store.GetObjectKey(i), err)
	}

	u.log.V(1).Info("issuer upserted", "resource", store.GetObjectKey(i), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertStatefulSet(ss *v1.StatefulSet, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert issuer", "resource", store.GetObjectKey(ss), "generation", gen)

	current := &v1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ss.GetName(),
			Namespace: ss.GetNamespace(),
		},
	}
	mgrClient := u.manager.GetClient()
	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		err := mergeMetadata(current, ss)
		if err != nil {
			return err
		}

		current.Spec.ServiceName = ss.Spec.ServiceName
		current.Spec.Selector = ss.Spec.Selector
		if ss.Spec.Replicas != nil {
			current.Spec.Replicas = ss.Spec.Replicas
		}
		//ss.Spec.Selector.DeepCopyInto(current.Spec.Selector)
		ss.Spec.Template.ObjectMeta.DeepCopyInto(&current.Spec.Template.ObjectMeta)

		currentSpec := &current.Spec.Template.Spec
		ssSpec := &ss.Spec.Template.Spec

		currentSpec.Volumes = make([]corev1.Volume, len(ssSpec.Volumes))
		for i := range ssSpec.Volumes {
			ssSpec.Volumes[i].DeepCopyInto(&currentSpec.Volumes[i])
		}

		currentSpec.Containers = make([]corev1.Container, len(ssSpec.Containers))
		for i := range ssSpec.Containers {
			ssSpec.Containers[i].DeepCopyInto(&currentSpec.Containers[i])
		}

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert statefulset %q: %w",
			store.GetObjectKey(ss), err)
	}

	u.log.V(1).Info("statefulset upserted", "resource", store.GetObjectKey(ss), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertGatewayClass(gwClass *gwapiv1.GatewayClass, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert gatewayclass", "resource", store.GetObjectKey(gwClass), "generation", gen)

	current := &gwapiv1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: gwClass.GetName(),
		},
	}
	mgrClient := u.manager.GetClient()
	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		err := mergeMetadata(current, gwClass)
		if err != nil {
			return err
		}

		gwClassSpec := &gwClass.Spec
		gwClassSpec.DeepCopyInto(&current.Spec)

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert gatewayclass %q: %w",
			store.GetObjectKey(gwClass), err)
	}

	u.log.V(1).Info("gatewayclass upserted", "resource", store.GetObjectKey(gwClass), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertGatewayConfigs(gwConfig *stnrgwv1.GatewayConfig, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert gatewayconfig", "resource", store.GetObjectKey(gwConfig), "generation", gen)

	current := &stnrgwv1.GatewayConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gwConfig.GetName(),
			Namespace: gwConfig.GetNamespace(),
		},
	}
	mgrClient := u.manager.GetClient()
	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		err := mergeMetadata(current, gwConfig)
		if err != nil {
			return err
		}

		gwConfigSpec := &gwConfig.Spec
		gwConfigSpec.DeepCopyInto(&current.Spec)

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert gatewayconfig %q: %w",
			store.GetObjectKey(gwConfig), err)
	}

	u.log.V(1).Info("gatewayconfig upserted", "resource", store.GetObjectKey(gwConfig), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertGateway(gw *gwapiv1.Gateway, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert gateway", "resource", store.GetObjectKey(gw), "generation", gen)

	current := &gwapiv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gw.GetName(),
			Namespace: gw.GetNamespace(),
		},
	}
	mgrClient := u.manager.GetClient()
	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		err := mergeMetadata(current, gw)
		if err != nil {
			return err
		}

		gwConfigSpec := &gw.Spec
		currentSpec := &current.Spec
		currentSpec.GatewayClassName = gwConfigSpec.GatewayClassName
		currentSpec.Listeners = make([]gwapiv1.Listener, len(gwConfigSpec.Listeners))
		for i, _ := range gwConfigSpec.Listeners {
			currentSpec.Listeners[i] = gwConfigSpec.Listeners[i]
		}

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert gateway %q: %w",
			store.GetObjectKey(gw), err)
	}

	u.log.V(1).Info("gateway upserted", "resource", store.GetObjectKey(gw), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertUDPRoute(udpr *stnrgwv1.UDPRoute, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert udproute", "resource", store.GetObjectKey(udpr), "generation", gen)

	current := &stnrgwv1.UDPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      udpr.GetName(),
			Namespace: udpr.GetNamespace(),
		},
	}
	mgrClient := u.manager.GetClient()
	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		err := mergeMetadata(current, udpr)
		if err != nil {
			return err
		}

		udprSpec := &udpr.Spec
		currentSpec := &current.Spec
		udprSpec.DeepCopyInto(currentSpec)

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert udproute %q: %w",
			store.GetObjectKey(udpr), err)
	}

	u.log.V(1).Info("udproute upserted", "resource", store.GetObjectKey(udpr), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertHTTPRoute(httpr *gwapiv1.HTTPRoute, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert httproute", "resource", store.GetObjectKey(httpr), "generation", gen)

	current := &gwapiv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      httpr.GetName(),
			Namespace: httpr.GetNamespace(),
		},
	}
	mgrClient := u.manager.GetClient()
	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		err := mergeMetadata(current, httpr)
		if err != nil {
			return err
		}

		httprSpec := &httpr.Spec
		currentSpec := &current.Spec
		httprSpec.DeepCopyInto(currentSpec)

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert httproute %q: %w",
			store.GetObjectKey(httpr), err)
	}

	u.log.V(1).Info("httpproute upserted", "resource", store.GetObjectKey(httpr), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertServiceAccount(sa *corev1.ServiceAccount, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert serviceaccount", "resource", store.GetObjectKey(sa), "generation", gen)

	current := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sa.GetName(),
			Namespace: sa.GetNamespace(),
		},
	}

	mgrClient := u.manager.GetClient()
	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		err := mergeMetadata(current, sa)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert serviceaccount %q: %w",
			store.GetObjectKey(sa), err)
	}

	u.log.V(1).Info("serviceaccount upserted", "resource", store.GetObjectKey(sa), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertClusterRole(r *rbacv1.ClusterRole, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert clusterrole", "resource", store.GetObjectKey(r), "generation", gen)

	current := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.GetName(),
		},
	}

	mgrClient := u.manager.GetClient()
	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		mergeMetadataForClusterScopedResources(current, r)

		current.Rules = make([]rbacv1.PolicyRule, len(r.Rules))
		for i := range r.Rules {
			r.Rules[i].DeepCopyInto(&current.Rules[i])
		}

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert clusterrole %q: %w",
			store.GetObjectKey(r), err)
	}

	u.log.V(1).Info("clusterrole upserted", "resource", store.GetObjectKey(r), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertClusterRoleBinding(rb *rbacv1.ClusterRoleBinding, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert clusterrolebinding", "resource", store.GetObjectKey(rb), "generation", gen)

	current := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: rb.GetName(),
		},
	}

	mgrClient := u.manager.GetClient()
	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		mergeMetadataForClusterScopedResources(current, rb)

		rb.RoleRef.DeepCopyInto(&current.RoleRef)

		current.Subjects = make([]rbacv1.Subject, len(rb.Subjects))
		for i := range rb.Subjects {
			rb.Subjects[i].DeepCopyInto(&current.Subjects[i])
		}

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert clusterrolebinding %q: %w",
			store.GetObjectKey(rb), err)
	}

	u.log.V(1).Info("clusterrolebinding upserted", "resource", store.GetObjectKey(rb), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

	return op, nil
}

func (u *Updater) upsertTCPRoute(tcpr *gwapiv1a2.TCPRoute, gen int) (ctrlutil.OperationResult, error) {
	u.log.V(2).Info("upsert httproute", "resource", store.GetObjectKey(tcpr), "generation", gen)

	current := &gwapiv1a2.TCPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tcpr.GetName(),
			Namespace: tcpr.GetNamespace(),
		},
	}
	mgrClient := u.manager.GetClient()
	op, err := ctrlutil.CreateOrUpdate(u.ctx, mgrClient, current, func() error {
		err := mergeMetadata(current, tcpr)
		if err != nil {
			return err
		}

		tcprSpec := &tcpr.Spec
		currentSpec := &current.Spec
		tcprSpec.DeepCopyInto(currentSpec)

		return nil
	})

	if err != nil {
		return ctrlutil.OperationResultNone, fmt.Errorf("cannot upsert tcproute %q: %w",
			store.GetObjectKey(tcpr), err)
	}

	u.log.V(1).Info("tcpproute upserted", "resource", store.GetObjectKey(tcpr), "generation",
		gen, "result", store.GetObjectKey(current)) //store.DumpObject(current))

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

func mergeMetadataForClusterScopedResources(dst, src client.Object) {
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
