package kubernetes

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewInClusterClient creates a new Kubernetes client using in-cluster configuration
func NewInClusterClient() (*kubernetes.Clientset, error) {
	// Use in-cluster config when running inside a Kubernetes pod
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting in-cluster config: %w", err)
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes client: %w", err)
	}

	return clientset, nil
}