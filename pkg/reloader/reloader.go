package reloader

import (
	"fmt"
	"os"
	"time"

	"github.com/kuoss/common/logger"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type reloader struct {
	started    bool
	controller *ReloadController
}

func getClientset() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func New(clientset *kubernetes.Clientset) (*reloader, error) {
	var err error
	if clientset == nil {
		clientset, err = getClientset()
		if err != nil {
			return nil, fmt.Errorf("getClientset err: %w", err)
		}
	}
	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)
	controller, err := NewReloadController(factory)
	if err != nil {
		return nil, fmt.Errorf("NewReloadController err: %w", err)
	}
	return &reloader{
		started:    false,
		controller: controller,
	}, nil
}

func (r *reloader) Start() error {
	if r.started {
		return fmt.Errorf("already started")
	}
	r.started = true
	go r.run()
	return nil
}

func (r *reloader) run() {
	stop := make(chan struct{})
	defer close(stop)
	err := r.controller.Run(stop)
	if err != nil {
		fmt.Printf("Run err: %s\n", err.Error())
		os.Exit(1)
	}
	logger.Infof("run end")
}
