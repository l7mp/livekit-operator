package renderer

import (
	"github.com/go-logr/logr"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
)

func RenderLiveKitMesh(mesh *lkstnv1a1.LiveKitMesh) {
	//TODO render each component in the livekitmesh
	logr.Logger{}.Info("Renderlivekitmesh", "lkname", mesh.Name)
}
