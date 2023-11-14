package renderer

import (
	"context"
	// "fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"

	// gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	// stunnerconfv1alpha1 "github.com/l7mp/stunner/pkg/apis/v1alpha1"

	"github.com/l7mp/livekit-operator/internal/event"
)

type Config struct {
	Scheme *runtime.Scheme
	//Manager manager.Manager
	Logger logr.Logger
}

type Renderer struct {
	ctx context.Context
	//mgr                  manager.Manager
	scheme               *runtime.Scheme
	gen                  int
	renderCh, operatorCh chan event.Event
	log, logger          logr.Logger
}

// NewRenderer creates a new Renderer
func NewRenderer(cfg Config) *Renderer {
	return &Renderer{
		scheme:   cfg.Scheme,
		renderCh: make(chan event.Event, 10),
		//mgr:      cfg.Manager,
		gen:    0,
		log:    cfg.Logger.WithName("renderer"),
		logger: cfg.Logger,
	}
}

func (r *Renderer) Start(ctx context.Context) error {
	r.ctx = ctx

	// starting the renderer thread
	go func() {
		defer close(r.renderCh)

		for {
			select {
			case e := <-r.renderCh:
				if e.GetType() != event.TypeRender {
					r.log.Info("renderer thread received unknown event",
						"event", e.GetType().String())
					continue
				}

				// prepare a new update event Render will populate
				// config is returned in the update event ConfigMap store
				ev := e.(*event.Render)
				r.RenderLiveKitMesh(ev)

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// GetRenderChannel returns the channel onn which the renderer listens to rendering requests
func (r *Renderer) GetRenderChannel() chan event.Event {
	return r.renderCh
}

// SetOperatorChannel sets the channel on which the operator event dispatcher listens
func (r *Renderer) SetOperatorChannel(ch chan event.Event) {
	r.operatorCh = ch
}
