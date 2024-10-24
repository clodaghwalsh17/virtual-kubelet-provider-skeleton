package customprovider

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	cpuCapacity    = "20"
	podCapacity    = "100"
	memoryCapacity = "10Gi"
)

func (p *CustomProvider) ConfigureNode(name string, node *v1.Node) {
	node.ObjectMeta.Name = name
	node.ObjectMeta.Labels["type"] = "virtual-kubelet"
	node.ObjectMeta.Labels["kubernetes.io/role"] = "agent"
	node.ObjectMeta.Labels["alpha.service-controller.kubernetes.io/exclude-balancer"] = "true"
	node.ObjectMeta.Labels["node.kubernetes.io/exclude-from-external-load-balancers"] = "true"

	taint := v1.Taint{
		Key:    "virtual-kubelet.io/provider",
		Value:  "custom-provider",
		Effect: v1.TaintEffectNoSchedule,
	}
	node.Spec.Taints = []v1.Taint{taint}

	node.Status.NodeInfo.OperatingSystem = "Linux"
	node.Status.NodeInfo.Architecture = "amd64"
	node.Status.Capacity = p.capacity()
	node.Status.Allocatable = p.capacity()
	node.Status.Conditions = p.nodeConditions()
	node.Status.Addresses = p.nodeAddresses()
	node.Status.DaemonEndpoints = p.nodeDaemonEndpoints()
}

func (p *CustomProvider) capacity() v1.ResourceList {
	resourceList := v1.ResourceList{
		v1.ResourceCPU:    resource.MustParse(cpuCapacity),
		v1.ResourceMemory: resource.MustParse(memoryCapacity),
		v1.ResourcePods:   resource.MustParse(podCapacity),
	}

	return resourceList
}

func (p *CustomProvider) nodeConditions() []v1.NodeCondition {
	return []v1.NodeCondition{
		{
			Type:               "Ready",
			Status:             v1.ConditionTrue,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletReady",
			Message:            "kubelet is ready.",
		},
		{
			Type:               "OutOfDisk",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasSufficientDisk",
			Message:            "kubelet has sufficient disk space available",
		},
		{
			Type:               "MemoryPressure",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasSufficientMemory",
			Message:            "kubelet has sufficient memory available",
		},
		{
			Type:               "DiskPressure",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasNoDiskPressure",
			Message:            "kubelet has no disk pressure",
		},
		{
			Type:               "NetworkUnavailable",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "RouteCreated",
			Message:            "RouteController created a route",
		},
	}
}

func (p *CustomProvider) nodeAddresses() []v1.NodeAddress {
	return []v1.NodeAddress{}
}

func (p *CustomProvider) nodeDaemonEndpoints() v1.NodeDaemonEndpoints {
	return v1.NodeDaemonEndpoints{
		KubeletEndpoint: v1.DaemonEndpoint{
			Port: 10250,
		},
	}
}
