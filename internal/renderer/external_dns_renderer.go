package renderer

import (
	"github.com/l7mp/livekit-operator/internal/store"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Renderer) renderExternalDNSCloudFlareDeployment(renderContext *RenderContext) {
	log := r.logger.WithName("renderExternalDNSCloudFlareDeployment")

	log.V(2).Info("trying to render External DNS Deployment")
	lkMesh := renderContext.liveKitMesh

	deployment := createExternalDNSCloudFlareDeployment(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, deployment, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(deployment))
		return
	}

	renderContext.update.UpsertQueue.Deployments.Upsert(deployment)
	log.V(2).Info("Upserted ExternalDNS Deployment into UpsertQueue", "deployment", store.GetObjectKey(deployment))
}

func (r *Renderer) renderExternalDNSCloudFlareServiceAccount(renderContext *RenderContext) {
	log := r.logger.WithName("renderExternalDNSCloudFlareServiceAccount")

	log.V(2).Info("trying to render External DNS ServiceAccount")
	lkMesh := renderContext.liveKitMesh

	serviceAccount := createExternalDNSServiceAccount(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, serviceAccount, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(serviceAccount))
		return
	}

	renderContext.update.UpsertQueue.ServiceAccounts.Upsert(serviceAccount)
	log.V(2).Info("Upserted ExternalDNS ServiceAccount into UpsertQueue", "serviceaccount", store.GetObjectKey(serviceAccount))
}

func (r *Renderer) renderExternalDNSCloudFlareDNSClusterRole(renderContext *RenderContext) {
	log := r.logger.WithName("renderExternalDNSCloudFlareDNSClusterRole")

	log.V(2).Info("trying to render External DNS ClusterRole")
	lkMesh := renderContext.liveKitMesh

	clusterRole := createExternalDNSClusterRole(lkMesh)
	//if err := controllerutil.SetOwnerReference(lkMesh, clusterRole, r.scheme); err != nil {
	//	log.Error(err, "cannot set owner reference", "owner",
	//		store.GetObjectKey(lkMesh), "reference",
	//		store.GetObjectKey(clusterRole))
	//	return
	//}

	renderContext.update.UpsertQueue.ClusterRoles.Upsert(clusterRole)
	log.V(2).Info("Upserted ExternalDNS ClusterRole into UpsertQueue", "clusterrole", store.GetObjectKey(clusterRole))
}

func (r *Renderer) renderExternalDNSCloudFlareDNSClusterRoleBinding(renderContext *RenderContext) {
	log := r.logger.WithName("renderExternalDNSCloudFlareDNSClusterRoleBinding")

	log.V(2).Info("trying to render External DNS ClusterRoleBinding")
	lkMesh := renderContext.liveKitMesh

	clusterRoleBinding := createExternalDNSClusterRoleBinding(lkMesh)
	//if err := controllerutil.SetOwnerReference(lkMesh, clusterRoleBinding, r.scheme); err != nil {
	//	log.Error(err, "cannot set owner reference", "owner",
	//		store.GetObjectKey(lkMesh), "reference",
	//		store.GetObjectKey(clusterRoleBinding))
	//	return
	//}

	renderContext.update.UpsertQueue.ClusterRoleBindings.Upsert(clusterRoleBinding)
	log.V(2).Info("Upserted ExternalDNS ClusterRoleBinding into UpsertQueue", "clusterrolebinding", store.GetObjectKey(clusterRoleBinding))
}
