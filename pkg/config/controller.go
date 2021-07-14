package config

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"istio.io/istio/pilot/pkg/config/memory"
	istiomodel "istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/adsc"
	istioconfig "istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
	"istio.io/pkg/log"
)

var (
	controllerLog = log.RegisterScope("config-controller", "config-controller debugging", 0)
	// We need serviceentry and virtualservice to generate the envoyfiters
	configCollection = collection.NewSchemasBuilder().
				MustAdd(collections.IstioNetworkingV1Alpha3Serviceentries).
				MustAdd(collections.IstioNetworkingV1Alpha3Virtualservices).
				MustAdd(collections.IstioNetworkingV1Alpha3Destinationrules).
				MustAdd(collections.IstioNetworkingV1Alpha3Envoyfilters).Build()
)

// Controller watches Istio config xDS server and notifies the listeners when config changes.
type Controller struct {
	configServerAddr string
	Store            istiomodel.ConfigStore
	controller       istiomodel.ConfigStoreCache
}

// NewController creates a new Controller instance based on the provided arguments.
func NewController(configServerAddr string) *Controller {
	store := memory.Make(configCollection)
	return &Controller{
		configServerAddr: configServerAddr,
		Store:            store,
		controller:       memory.NewController(store),
	}
}

// Run until a signal is received, this function won't block
func (c *Controller) Run(stop <-chan struct{}) {
	go func() {
		for {
			xdsMCP, err := adsc.New(c.configServerAddr, &adsc.Config{
				Meta: istiomodel.NodeMetadata{
					Generator: "api",
				}.ToStruct(),
				InitialDiscoveryRequests: c.configInitialRequests(),
				BackoffPolicy:            backoff.NewConstantBackOff(time.Second),
			})
			if err != nil {
				controllerLog.Errorf("failed to dial XDS %s %v", c.configServerAddr, err)
				time.Sleep(5 * time.Second)
				continue
			}
			xdsMCP.Store = istiomodel.MakeIstioStore(c.controller)
			if err = xdsMCP.Run(); err != nil {
				controllerLog.Errorf("adsc: failed running %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			c.controller.Run(stop)
			return
		}
	}()
}

func (c *Controller) configInitialRequests() []*discovery.DiscoveryRequest {
	schemas := configCollection.All()
	requests := make([]*discovery.DiscoveryRequest, len(schemas))
	for i, schema := range schemas {
		requests[i] = &discovery.DiscoveryRequest{
			TypeUrl: schema.Resource().GroupVersionKind().String(),
		}
	}
	return requests
}

// RegisterEventHandler adds a handler to receive config update events for a configuration type
func (c *Controller) RegisterEventHandler() {
	handlerWrapper := func(prev istioconfig.Config, curr istioconfig.Config, event istiomodel.Event) {
		controllerLog.Infof("receive istio event %s, curr %s", event, fmt.Sprintf("%s/%s", curr.Namespace, curr.Name))
	}

	schemas := configCollection.All()
	for _, schema := range schemas {
		c.controller.RegisterEventHandler(schema.Resource().GroupVersionKind(), handlerWrapper)
	}
}
