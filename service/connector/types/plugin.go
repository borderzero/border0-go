package types

const (
	// PluginTypeAwsEc2Discovery is the plugin type for aws ec2 instance discovery.
	PluginTypeAwsEc2Discovery = "aws_ec2_discovery"

	// PluginTypeAwsEcsDiscovery is the plugin type for aws ecs cluster discovery.
	PluginTypeAwsEcsDiscovery = "aws_ecs_discovery"

	// PluginTypeAwsRdsDiscovery is the plugin type for aws rds db instance discovery.
	PluginTypeAwsRdsDiscovery = "aws_rds_discovery"

	// PluginTypeAwsSsmDiscovery is the plugin type for aws ssm target discovery.
	PluginTypeAwsSsmDiscovery = "aws_ssm_discovery"

	// PluginTypeDockerDiscovery is the plugin type for docker container discovery.
	PluginTypeDockerDiscovery = "docker_discovery"

	// PluginTypeKubernetesDiscovery is the plugin type for kubernetes pod discovery.
	PluginTypeKubernetesDiscovery = "kubernetes_discovery"

	// PluginTypeNetworkDiscovery is the plugin type for network service discovery.
	PluginTypeNetworkDiscovery = "network_discovery"
)

// PluginConfiguration represents configuration for a Border0 connector plugin.
type PluginConfiguration struct {
	PluginType string `json:"plugin_type"`
	Enabled    bool   `json:"enabled"`

	AwsEc2DiscoveryPluginConfiguration     *AwsEc2DiscoveryPluginConfiguration     `json:"aws_ec2_discovery_plugin_configuration,omitempty"`
	AwsEcsDiscoveryPluginConfiguration     *AwsEcsDiscoveryPluginConfiguration     `json:"aws_ecs_discovery_plugin_configuration,omitempty"`
	AwsRdsDiscoveryPluginConfiguration     *AwsRdsDiscoveryPluginConfiguration     `json:"aws_rds_discovery_plugin_configuration,omitempty"`
	AwsSsmDiscoveryPluginConfiguration     *AwsSsmDiscoveryPluginConfiguration     `json:"aws_ssm_discovery_plugin_configuration,omitempty"`
	DockerDiscoveryPluginConfiguration     *DockerDiscoveryPluginConfiguration     `json:"docker_discovery_plugin_configuration,omitempty"`
	KubernetesDiscoveryPluginConfiguration *KubernetesDiscoveryPluginConfiguration `json:"kubernetes_discovery_plugin_configuration,omitempty"`
	NetworkDiscoveryPluginConfiguration    *NetworkDiscoveryPluginConfiguration    `json:"local_network_discovery_plugin_configuration,omitempty"`
}

// AwsCredentials represents credentials and coonfiguration for authenticating against AWS APIs.
type AwsCredentials struct {
	AwsProfile         *string `json:"aws_profile,omitempty"`
	AwsAccessKeyId     *string `json:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey *string `json:"aws_secret_access_key,omitempty"`
	AwsSessionToken    *string `json:"aws_session_token,omitempty"`
}

// BaseAwsPluginConfiguration represents configuration fields shared across all AWS related plugins.
type BaseAwsPluginConfiguration struct {
	AwsCredentials *AwsCredentials `json:"aws_credentials,omitempty"`
	AwsRegions     []string        `json:"aws_regions"`
}

// BaseDiscoveryPluginConfiguration represents configuration fields shared across all discovery related plugins.
type BaseDiscoveryPluginConfiguration struct {
	ScanIntervalSeconds uint32 `json:"scan_interval_seconds"`
}

// AwsEc2DiscoveryPluginConfiguration represents configuration for the aws_ec2_discovery plugin.
type AwsEc2DiscoveryPluginConfiguration struct {
	BaseAwsPluginConfiguration       // extends
	BaseDiscoveryPluginConfiguration // extends

	IncludeWithStates []string            `json:"include_with_states,omitempty"`
	IncludeWithTags   map[string][]string `json:"include_with_tags,omitempty"`
	ExcludeWithTags   map[string][]string `json:"exclude_with_tags,omitempty"`
}

// AwsEcsDiscoveryPluginConfiguration represents configuration for the aws_ecs_discovery plugin.
type AwsEcsDiscoveryPluginConfiguration struct {
	BaseAwsPluginConfiguration       // extends
	BaseDiscoveryPluginConfiguration // extends

	IncludeWithStatuses []string            `json:"include_with_statuses,omitempty"`
	IncludeWithTags     map[string][]string `json:"include_with_tags,omitempty"`
	ExcludeWithTags     map[string][]string `json:"exclude_with_tags,omitempty"`
}

// AwsRdsDiscoveryPluginConfiguration represents configuration for the aws_rds_discovery plugin.
type AwsRdsDiscoveryPluginConfiguration struct {
	BaseAwsPluginConfiguration       // extends
	BaseDiscoveryPluginConfiguration // extends

	IncludeWithStatuses []string            `json:"include_with_statuses,omitempty"`
	IncludeWithTags     map[string][]string `json:"include_with_tags,omitempty"`
	ExcludeWithTags     map[string][]string `json:"exclude_with_tags,omitempty"`
}

// AwsSsmDiscoveryPluginConfiguration represents configuration for the aws_ssm_discovery plugin.
type AwsSsmDiscoveryPluginConfiguration struct {
	BaseAwsPluginConfiguration       // extends
	BaseDiscoveryPluginConfiguration // extends

	IncludeWithPingStatuses []string `json:"include_with_ping_statuses,omitempty"`
}

// DockerDiscoveryPluginConfiguration represents configuration for the docker_discovery plugin.
type DockerDiscoveryPluginConfiguration struct {
	BaseDiscoveryPluginConfiguration // extends

	// TODO: docker host unix socket, etc...
}

// KubernetesDiscoveryPluginConfiguration represents configuration for the kubernetes_discovery plugin.
type KubernetesDiscoveryPluginConfiguration struct {
	BaseDiscoveryPluginConfiguration // extends

	// TODO: k8s api URL, credentials...
}

// NetworkDiscoveryPluginConfiguration represents configuration for the network_discovery plugin.
type NetworkDiscoveryPluginConfiguration struct {
	BaseDiscoveryPluginConfiguration // extends

	// map key can be a hostname, an IP, IP-range, or CIDR
	Targets map[string]struct {
		Ports                      []uint16 `json:"ports"`
		DiscoverySshServers        bool     `json:"discover_ssh_servers,omitempty"`
		DiscoveryHttpServers       bool     `json:"discover_http_servers,omitempty"`
		DiscoveryMysqlServers      bool     `json:"discover_mysql_servers,omitempty"`
		DiscoveryPostgresqlServers bool     `json:"discover_postgresql_servers,omitempty"`
	} `json:"targets"`
}
