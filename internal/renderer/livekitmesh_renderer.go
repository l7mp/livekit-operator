package renderer

import (
	"github.com/l7mp/livekit-operator/internal/event"
	"github.com/l7mp/livekit-operator/internal/store"
)

func (r *Renderer) RenderLiveKitMesh(e *event.Render) {
	log := r.logger.WithName("RenderLiveKitMesh")
	log.Info("Trying to render LiveKitMeshes")
	//TODO render each component in the livekitmesh
	r.gen += 1

	if liveKitMeshes := store.LiveKitMeshes.GetAll(); len(liveKitMeshes) == 0 {
		log.Info("no LiveKitMesh objects found", "event", e.String())
		//TODO maybe return here?
	} else {

		for _, lkMesh := range liveKitMeshes {
			log.Info("Found in store", "lk", lkMesh.Name)
			lkMesh := lkMesh
			renderContext := NewRenderContext(e, r, lkMesh)
			_ = renderContext
			//renderContext.lkMesh
		}
	}
}
