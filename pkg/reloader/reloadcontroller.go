package reloader

import (
	"fmt"
	"os"

	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type ReloadController struct {
	informerFactory informers.SharedInformerFactory
	serviceInformer coreinformers.ServiceInformer
}

func NewReloadController(informerFactory informers.SharedInformerFactory) (*ReloadController, error) {
	serviceInformer := informerFactory.Core().V1().Services()

	c := &ReloadController{
		informerFactory: informerFactory,
		serviceInformer: serviceInformer,
	}
	_, err := serviceInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				c.Stop()
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				c.Stop()
			},
			DeleteFunc: func(obj interface{}) {
				c.Stop()
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("AddEventHandler err: %w", err)
	}
	return c, nil
}

func (c *ReloadController) Run(ch chan struct{}) error {
	c.informerFactory.Start(ch)
	if !cache.WaitForCacheSync(ch, c.serviceInformer.Informer().HasSynced) {
		return fmt.Errorf("cannot sync")
	}
	return nil
}

func (c *ReloadController) Stop() {
	os.Exit(0)
}
