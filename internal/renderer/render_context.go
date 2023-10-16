package renderer

import (
	"github.com/go-logr/logr"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/event"
)

// RenderContext contains all the components from the LiveKitMesh configuration for the current rendering task,
type RenderContext struct {
	origin      event.Event
	update      *event.Update
	lkMesh      *lkstnv1a1.LiveKit
	ingress     *lkstnv1a1.Ingress
	egress      *lkstnv1a1.Egress
	certManager *lkstnv1a1.CertManager
	monitoring  *lkstnv1a1.Monitoring
	//scaling
	log logr.Logger
}

func NewRenderContext(e *event.Render, r *Renderer, lkMesh *lkstnv1a1.LiveKitMesh) *RenderContext {
	return &RenderContext{
		origin:      e,
		update:      event.NewEventUpdate(r.gen),
		lkMesh:      lkMesh.Spec.Components.LiveKit,
		ingress:     lkMesh.Spec.Components.Ingress,
		egress:      lkMesh.Spec.Components.Egress,
		certManager: lkMesh.Spec.Components.CertManager,
		monitoring:  lkMesh.Spec.Components.Monitoring,
		log:         r.log.WithValues("LiveKitMesh", lkMesh.GetName()),
	}
}
