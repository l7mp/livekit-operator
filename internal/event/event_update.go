package event

import (
	"github.com/l7mp/livekit-operator/internal/store"
)

// UpdateConf render event
type UpdateConf struct {
	ConfigMaps    *store.ConfigMapStore
	LiveKitMeshes *store.LiveKitMeshStore
	Services      *store.ServiceStore
	Deployments   *store.DeploymentStore
	Issuer        *store.IssuerStore
	Secret        *store.SecretStore
}

type Update struct {
	Type        Type
	UpsertQueue UpdateConf
	DeleteQueue UpdateConf
	Generation  int
}

// NewEventUpdate returns an empty event
func NewEventUpdate(generation int) *Update {
	return &Update{
		Type: TypeUpdate,
		UpsertQueue: UpdateConf{
			LiveKitMeshes: store.NewLivekitMeshStore(),
			ConfigMaps:    store.NewConfigMapStore(),
			Services:      store.NewServiceStore(),
			Deployments:   store.NewDeploymentStore(),
			Issuer:        store.NewIssuerStore(),
			Secret:        store.NewSecretStore(),
		},
		DeleteQueue: UpdateConf{
			LiveKitMeshes: store.NewLivekitMeshStore(),
			ConfigMaps:    store.NewConfigMapStore(),
			//FIXME not everything should be deleted/handled by us that what other operators should do
		},
		Generation: generation,
	}
}

func (e *Update) GetType() Type {
	return e.Type
}

func (e *Update) String() string {
	//return fmt.Sprintf("%s (gen: %d): upsert-queue: gway-cls: %d, gway: %d, route: %d, svc: %d, confmap: %d, dp: %d / "+
	//	"delete-queue: gway-cls: %d, gway: %d, route: %d, svc: %d, confmap: %d, dp: %d", e.Type.String(),
	//	e.Generation, e.UpsertQueue.GatewayClasses.Len(), e.UpsertQueue.Gateways.Len(),
	//	e.UpsertQueue.UDPRoutes.Len(), e.UpsertQueue.Services.Len(),
	//	e.UpsertQueue.ConfigMaps.Len(), e.UpsertQueue.Deployments.Len(),
	//	e.DeleteQueue.GatewayClasses.Len(), e.DeleteQueue.Gateways.Len(),
	//	e.DeleteQueue.UDPRoutes.Len(), e.DeleteQueue.Services.Len(),
	//	e.DeleteQueue.ConfigMaps.Len(), e.DeleteQueue.Deployments.Len())
	return "fixme"
	//TODO
}
