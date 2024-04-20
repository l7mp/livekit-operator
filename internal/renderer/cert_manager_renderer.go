package renderer

import (
	"github.com/l7mp/livekit-operator/internal/store"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Renderer) renderCertManagerIssuerAndSecret(context *RenderContext) {
	log := r.logger.WithName("renderCertManagerIssuer")

	log.V(2).Info("trying to render Issuer", "component", "Cert-Manager")
	lkMesh := context.liveKitMesh

	issuer, secret := createIssuer(*lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, issuer, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(issuer))
		return
	}
	if err := controllerutil.SetOwnerReference(lkMesh, secret, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(secret))
		return
	}

	context.update.UpsertQueue.Issuer.Upsert(issuer)

	log.V(2).Info("Upserted Cert-Manager Issuer into UpsertQueue", "cm", store.GetObjectKey(issuer))

	context.update.UpsertQueue.Secrets.Upsert(secret)

	log.V(2).Info("Upserted Cert-Manager Secrets into UpsertQueue", "cm", store.GetObjectKey(secret))
}
