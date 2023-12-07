package connector

import "github.com/borderzero/border0-go/types/common"

const (
	// PluginTypeAwsEc2Discovery is the plugin type for aws ec2 instance discovery.
	PluginTypeAwsEc2Discovery = "aws_ec2_discovery"

	// PluginTypeAwsEcsDiscovery is the plugin type for aws ecs service discovery.
	PluginTypeAwsEcsDiscovery = "aws_ecs_discovery"

	// PluginTypeAwsEksDiscovery is the plugin type for aws eks cluster discovery.
	PluginTypeAwsEksDiscovery = "aws_eks_discovery"

	// PluginTypeAwsRdsDiscovery is the plugin type for aws rds db instance discovery.
	PluginTypeAwsRdsDiscovery = "aws_rds_discovery"

	// PluginTypeDockerDiscovery is the plugin type for docker container discovery.
	PluginTypeDockerDiscovery = "docker_discovery"

	// PluginTypeKubernetesDiscovery is the plugin type for kubernetes pod discovery.
	PluginTypeKubernetesDiscovery = "kubernetes_discovery"

	// PluginTypeNetworkDiscovery is the plugin type for network service discovery.
	PluginTypeNetworkDiscovery = "network_discovery"
)

// PluginConfiguration represents configuration for a Border0 connector plugin.
type PluginConfiguration struct {
	AwsEc2DiscoveryPluginConfiguration     *AwsEc2DiscoveryPluginConfiguration     `json:"aws_ec2_discovery_plugin_configuration,omitempty"`
	AwsEcsDiscoveryPluginConfiguration     *AwsEcsDiscoveryPluginConfiguration     `json:"aws_ecs_discovery_plugin_configuration,omitempty"`
	AwsRdsDiscoveryPluginConfiguration     *AwsRdsDiscoveryPluginConfiguration     `json:"aws_rds_discovery_plugin_configuration,omitempty"`
	AwsEksDiscoveryPluginConfiguration     *AwsEksDiscoveryPluginConfiguration     `json:"aws_eks_discovery_plugin_configuration,omitempty"`
	DockerDiscoveryPluginConfiguration     *DockerDiscoveryPluginConfiguration     `json:"docker_discovery_plugin_configuration,omitempty"`
	KubernetesDiscoveryPluginConfiguration *KubernetesDiscoveryPluginConfiguration `json:"kubernetes_discovery_plugin_configuration,omitempty"`
	NetworkDiscoveryPluginConfiguration    *NetworkDiscoveryPluginConfiguration    `json:"network_discovery_plugin_configuration,omitempty"`
}

// KubernetesCredentials represents credentials and configuration for authenticating against a Kubernetes API.
type KubernetesCredentials struct {
	MasterUrl      *string `json:"master_url,omitempty"`
	KubeconfigPath *string `json:"kubeconfig_path,omitempty"`
}

// BaseAwsPluginConfiguration represents configuration fields shared across all AWS related plugins.
type BaseAwsPluginConfiguration struct {
	AwsCredentials *common.AwsCredentials `json:"aws_credentials,omitempty"`
	AwsRegions     []string               `json:"aws_regions,omitempty"`
}

// BaseDiscoveryPluginConfiguration represents configuration fields shared across all discovery related plugins.
type BaseDiscoveryPluginConfiguration struct {
	ScanIntervalMinutes uint32 `json:"scan_interval_minutes"`
}

// AwsEc2DiscoveryPluginConfiguration represents configuration for the aws_ec2_discovery plugin.
type AwsEc2DiscoveryPluginConfiguration struct {
	BaseAwsPluginConfiguration       // extends
	BaseDiscoveryPluginConfiguration // extends

	IncludeWithStates []string            `json:"include_with_states,omitempty"`
	IncludeWithTags   map[string][]string `json:"include_with_tags,omitempty"`
	ExcludeWithTags   map[string][]string `json:"exclude_with_tags,omitempty"`
	CheckSsmStatus    bool                `json:"check_ssm_status"`
}

// AwsEcsDiscoveryPluginConfiguration represents configuration for the aws_ecs_discovery plugin.
type AwsEcsDiscoveryPluginConfiguration struct {
	BaseAwsPluginConfiguration       // extends
	BaseDiscoveryPluginConfiguration // extends

	IncludeWithTags map[string][]string `json:"include_with_tags,omitempty"`
	ExcludeWithTags map[string][]string `json:"exclude_with_tags,omitempty"`
}

// AwsEksDiscoveryPluginConfiguration represents configuration for the aws_eks_discovery plugin.
type AwsEksDiscoveryPluginConfiguration struct {
	BaseAwsPluginConfiguration       // extends
	BaseDiscoveryPluginConfiguration // extends

	IncludeWithTags map[string][]string `json:"include_with_tags,omitempty"`
	ExcludeWithTags map[string][]string `json:"exclude_with_tags,omitempty"`
}

// AwsRdsDiscoveryPluginConfiguration represents configuration for the aws_rds_discovery plugin.
type AwsRdsDiscoveryPluginConfiguration struct {
	BaseAwsPluginConfiguration       // extends
	BaseDiscoveryPluginConfiguration // extends

	IncludeWithStatuses []string            `json:"include_with_statuses,omitempty"`
	IncludeWithTags     map[string][]string `json:"include_with_tags,omitempty"`
	ExcludeWithTags     map[string][]string `json:"exclude_with_tags,omitempty"`
}

// DockerDiscoveryPluginConfiguration represents configuration for the docker_discovery plugin.
type DockerDiscoveryPluginConfiguration struct {
	BaseDiscoveryPluginConfiguration // extends

	IncludeWithLabels map[string][]string `json:"include_with_labels,omitempty"`
	ExcludeWithLabels map[string][]string `json:"exclude_with_labels,omitempty"`
}

// KubernetesDiscoveryPluginConfiguration represents configuration for the kubernetes_discovery plugin.
type KubernetesDiscoveryPluginConfiguration struct {
	BaseDiscoveryPluginConfiguration // extends

	KubernetesCredentials *KubernetesCredentials `json:"kubernetes_credentials,omitempty"`

	Namespaces []string `json:"namespaces,omitempty"`

	IncludeWithLabels map[string][]string `json:"include_with_labels,omitempty"`
	ExcludeWithLabels map[string][]string `json:"exclude_with_labels,omitempty"`
}

// NetworkDiscoveryTarget represents a single target and configuration for the network_discovery plugin.
type NetworkDiscoveryTarget struct {
	Target string   `json:"target"`
	Ports  []uint16 `json:"ports"`
}

// NetworkDiscoveryPluginConfiguration represents configuration for the network_discovery plugin.
type NetworkDiscoveryPluginConfiguration struct {
	BaseDiscoveryPluginConfiguration // extends

	Targets []NetworkDiscoveryTarget `json:"targets"`
}
