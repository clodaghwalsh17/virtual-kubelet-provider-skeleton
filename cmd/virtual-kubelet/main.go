package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/virtual-kubelet/virtual-kubelet/node"
	"github.com/virtual-kubelet/virtual-kubelet/node/nodeutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	customProvider "customVK/pkg/provider"
)

var (
	nodeName string
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	desc := "A description of the tool"

	cmd := &cobra.Command{
		Use:   "virtual-kubelet",
		Short: desc,
		Long:  desc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(ctx)
		},
	}

	cmd.Flags().StringVar(&nodeName, "nodename", "vk-node", "Name of the node")

	if err := cmd.Execute(); err != nil {
		fmt.Println("error running the virtual-kubelet command: %w", err)
	}
}

func run(ctx context.Context) error {
	fmt.Println(nodeName)

	provider := func(cfg nodeutil.ProviderConfig) (nodeutil.Provider, node.NodeProvider, error) {
		p, err := customProvider.NewCustomProvider(nodeName, cfg)

		if err != nil {
			return nil, nil, err
		}

		return p, nil, nil
	}

	// Set K8s Client for Provider

	// Use for getting cluster config when the code is ran inside a cluster - when pkg and run as K8s deployment
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	k8sClient, err := kubernetes.NewForConfig(config)

	// Create client with kubeconfig file use if running the provider as a go binary
	// kubeConfigPath := "/home/.kube/config"
	// k8sClient, err := nodeutil.ClientsetFromEnv(kubeConfigPath)

	if err != nil {
		panic(err.Error())
	}

	withClient := func(cfg *nodeutil.NodeConfig) error {
		return nodeutil.WithClient(k8sClient)(cfg)
	}

	attachRoutes := func(cfg *nodeutil.NodeConfig) error {
		mux := http.NewServeMux()
		cfg.Handler = mux
		return nodeutil.AttachProviderRoutes(mux)(cfg)
	}

	additionalConfig := func(cfg *nodeutil.NodeConfig) error {
		cfg.HTTPListenAddr = ":10250"
		cfg.NumWorkers = 50
		cfg.DebugHTTP = true
		return nil
	}

	node, err := nodeutil.NewNode(nodeName, provider, withClient, attachRoutes, additionalConfig)

	if err != nil {
		return err
	}

	go func() error {
		err = node.Run(ctx)
		if err != nil {
			return fmt.Errorf("error running the node: %w", err)
		}
		return nil
	}()

	<-node.Done()
	return node.Err()
}
