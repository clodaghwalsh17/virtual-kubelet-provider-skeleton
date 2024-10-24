package customprovider

import (
	"context"
	"fmt"
	"io"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	statsv1alpha1 "github.com/virtual-kubelet/virtual-kubelet/node/api/statsv1alpha1"
	"github.com/virtual-kubelet/virtual-kubelet/node/nodeutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CustomProvider struct {
	name string
}

func NewCustomProvider(name string, cfg nodeutil.ProviderConfig) (*CustomProvider, error) {
	p := CustomProvider{}
	p.name = name

	p.ConfigureNode(name, cfg.Node)

	return &p, nil
}

// All errors for CustomProvider functions implemented need to be caught using github.com/virtual-kubelet/virtual-kubelet/errdefs

func (p *CustomProvider) CreatePod(ctx context.Context, pod *v1.Pod) error {
	fmt.Println("Called created fn")

	now := metav1.NewTime(time.Now())
	pod.Status = v1.PodStatus{
		Phase:     v1.PodRunning,
		HostIP:    "1.2.3.4",
		PodIP:     "5.6.7.8",
		StartTime: &now,
		Conditions: []v1.PodCondition{
			{
				Type:   v1.PodInitialized,
				Status: v1.ConditionTrue,
			},
			{
				Type:   v1.PodReady,
				Status: v1.ConditionTrue,
			},
			{
				Type:   v1.PodScheduled,
				Status: v1.ConditionTrue,
			},
		},
	}

	return nil
}

func (p *CustomProvider) UpdatePod(ctx context.Context, pod *v1.Pod) error {
	fmt.Println("Called update fn")
	return nil
}

func (p *CustomProvider) DeletePod(ctx context.Context, pod *v1.Pod) error {
	fmt.Println("Called delete fn")
	return nil
}

func (p *CustomProvider) GetPod(ctx context.Context, namespace, name string) (*v1.Pod, error) {
	fmt.Println("Called get pod fn")
	return nil, nil
}

func (p *CustomProvider) GetPodStatus(ctx context.Context, namespace, name string) (*v1.PodStatus, error) {
	fmt.Println("Called get pod status fn")
	return &v1.PodStatus{}, nil
}

func (p *CustomProvider) GetPods(context.Context) ([]*v1.Pod, error) {
	fmt.Println("Called get pods fn")
	return nil, nil
}

func (p *CustomProvider) GetContainerLogs(ctx context.Context, namespace, podName, containerName string, opts api.ContainerLogOpts) (io.ReadCloser, error) {
	fmt.Println("Called get container logs fn")
	return nil, nil
}

func (p *CustomProvider) RunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string, attach api.AttachIO) error {
	fmt.Println("Called run in container fn")
	return nil
}

func (p *CustomProvider) AttachToContainer(ctx context.Context, namespace, podName, containerName string, attach api.AttachIO) error {
	fmt.Println("Called attach to container fn")
	return nil
}

func (p *CustomProvider) PortForward(ctx context.Context, namespace, pod string, port int32, stream io.ReadWriteCloser) error {
	fmt.Println("Called port forward fn")
	return nil
}

func (p *CustomProvider) GetStatsSummary(context.Context) (*statsv1alpha1.Summary, error) {
	fmt.Println("Called get stats summary fn")
	return &statsv1alpha1.Summary{}, nil
}

func (p *CustomProvider) GetMetricsResource(context.Context) ([]*dto.MetricFamily, error) {
	fmt.Println("Called get metrics resource fn")
	return []*dto.MetricFamily{}, nil
}
