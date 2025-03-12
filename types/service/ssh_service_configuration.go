package service

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/borderzero/border0-go/lib/types/nilcheck"
	"github.com/borderzero/border0-go/lib/types/set"
	"github.com/borderzero/border0-go/types/common"
	"golang.org/x/crypto/ssh"
)

const (
	// SshServiceTypeStandard is the ssh
	// service type for standard ssh services.
	SshServiceTypeStandard = "standard"

	// SshServiceTypeAwsSsm is the ssh service
	// type for aws session manager ssh services.
	SshServiceTypeAwsSsm = "aws_ssm"

	// SshServiceTypeAwsEc2InstanceConnect is the ssh service
	// type for aws ec2 instance connect ssh services.
	SshServiceTypeAwsEc2InstanceConnect = "aws_ec2_instance_connect"

	// SshServiceTypeKubectlExec is the ssh service
	// type for kubectl exec ssh services.
	SshServiceTypeKubectlExec = "kubectl_exec"

	// SshServiceTypeDockerExec is the ssh service
	// type for docker exec ssh services.
	SshServiceTypeDockerExec = "docker_exec"

	// SshServiceTypeConnectorBuiltIn is the ssh service
	// type for the connector's built-in ssh service.
	SshServiceTypeConnectorBuiltIn = "connector_built_in_ssh_service"
)

const (
	// SsmTargetTypeEc2 is the ssm target type for ec2 targets.
	SsmTargetTypeEc2 = "ec2"

	// SsmTargetTypeEcs is the ssm target type for ecs targets.
	SsmTargetTypeEcs = "ecs"
)

const (
	// KubectlExecTargetTypeStandard is the kubectl
	// exec target type for standard k8s clusters.
	KubectlExecTargetTypeStandard = "standard"

	// KubectlExecTargetTypeAwsEks is the kubectl
	// exec target type for aws eks k8s clusters.
	KubectlExecTargetTypeAwsEks = "aws_eks"
)

const (
	// StandardSshServiceAuthenticationTypeUsernameAndPassword is the standard ssh
	// service authentication type for authenticating with a username and password.
	StandardSshServiceAuthenticationTypeUsernameAndPassword = "username_and_password"

	// StandardSshServiceAuthenticationTypePrivateKey is the standard ssh
	// service authentication type for authenticating with a private key.
	StandardSshServiceAuthenticationTypePrivateKey = "private_key"

	// StandardSshServiceAuthenticationTypeBorder0Certificate is the standard ssh
	// service authentication type for authenticating with a border0-signed certificate.
	StandardSshServiceAuthenticationTypeBorder0Certificate = "border0_certificate"
)

const (
	// UsernameProviderDefined is the username provider
	// option for using an admin-defined (static) username.
	UsernameProviderDefined = "defined"

	// UsernameProviderPromptClient is username provider option
	// for prompting connecting clients for the username.
	UsernameProviderPromptClient = "prompt_client"

	// UsernameProviderUseConnectorUser is username provider
	// option for using the connector's OS username.
	//
	// NOTE: This option can only be used as the username
	// provider for connector built-in ssh services.
	UsernameProviderUseConnectorUser = "use_connector_user"
)

const (
	ec2InstanceConnectEndpointIdRegex = `^eice-[0-9a-f]{17}$`
)

// SshServiceConfiguration represents service
// configuration for shell services (fka sockets).
type SshServiceConfiguration struct {
	SshServiceType string `json:"ssh_service_type"`

	// mutually exclusive fields below
	StandardSshServiceConfiguration    *StandardSshServiceConfiguration    `json:"standard_ssh_service_configuration,omitempty"`
	AwsSsmSshServiceConfiguration      *AwsSsmSshServiceConfiguration      `json:"aws_ssm_ssh_service_configuration,omitempty"`
	AwsEc2ICSshServiceConfiguration    *AwsEc2ICSshServiceConfiguration    `json:"aws_ec2ic_ssh_service_configuration,omitempty"`
	DockerExecSshServiceConfiguration  *DockerExecSshServiceConfiguration  `json:"docker_exec_ssh_service_configuration,omitempty"`
	KubectlExecSshServiceConfiguration *KubectlExecSshServiceConfiguration `json:"kubectl_exec_ssh_service_configuration,omitempty"`
	BuiltInSshServiceConfiguration     *BuiltInSshServiceConfiguration     `json:"built_in_ssh_service_configuration,omitempty"`
}

// StandardSshServiceConfiguration represents service
// configuration for standard ssh services (fka sockets).
type StandardSshServiceConfiguration struct {
	HostnameAndPort
	SshAuthenticationType string `json:"ssh_authentication_type"`

	// mutually exclusive fields below
	UsernameAndPasswordAuthConfiguration *UsernameAndPasswordAuthConfiguration `json:"username_and_password_auth_configuration,omitempty"`
	PrivateKeyAuthConfiguration          *PrivateKeyAuthConfiguration          `json:"private_key_auth_configuration,omitempty"`
	Border0CertificateAuthConfiguration  *Border0CertificateAuthConfiguration  `json:"border0_certificate_auth_configuration,omitempty"`
}

// AwsSsmSshServiceConfiguration represents service
// configuration for aws ssm ssh services (fka sockets).
type AwsSsmSshServiceConfiguration struct {
	SsmTargetType string `json:"ssm_target_type"`

	// mutually exclusive fields below
	AwsSsmEc2TargetConfiguration *AwsSsmEc2TargetConfiguration `json:"aws_ssm_ec2_target_configuration,omitempty"`
	AwsSsmEcsTargetConfiguration *AwsSsmEcsTargetConfiguration `json:"aws_ssm_ecs_target_configuration,omitempty"`
}

// AwsSsmEc2TargetConfiguration represents service configuration
// for aws ssm ssh services (fka sockets) that have EC2 instances
// as their ssm target.
type AwsSsmEc2TargetConfiguration struct {
	Ec2InstanceId     string                 `json:"ec2_instance_id"`
	Ec2InstanceRegion string                 `json:"ec2_instance_region"`
	AwsCredentials    *common.AwsCredentials `json:"aws_credentials,omitempty"`
}

// AwsSsmEcsTargetConfiguration represents service configuration
// for aws ssm ssh services (fka sockets) that have ECS services
// as their ssm target.
type AwsSsmEcsTargetConfiguration struct {
	EcsClusterRegion string                 `json:"ecs_cluster_region"`
	EcsClusterName   string                 `json:"ecs_cluster_name"`
	EcsServiceName   string                 `json:"ecs_service_name"`
	AwsCredentials   *common.AwsCredentials `json:"aws_credentials,omitempty"`
}

// AwsEc2ICSshServiceConfiguration represents service configuration
// for aws ec2 instance connect ssh services (fka sockets).
type AwsEc2ICSshServiceConfiguration struct {
	HostnameAndPort
	UsernameProvider             string                 `json:"username_provider,omitempty"`
	Username                     string                 `json:"username,omitempty"`
	Ec2InstanceId                string                 `json:"ec2_instance_id"`
	Ec2InstanceRegion            string                 `json:"ec2_instance_region"`
	Ec2InstanceConnectEndpointId string                 `json:"ec2_instance_connect_endpoint_id,omitempty"`
	AwsCredentials               *common.AwsCredentials `json:"aws_credentials,omitempty"`
}

// DockerExecSshServiceConfiguration represents service
// configuration for docker exec ssh services (fka sockets).
type DockerExecSshServiceConfiguration struct {
	ContainerNameAllowlist []string `json:"container_name_allowlist,omitempty"`
}

// KubectlExecSshServiceConfiguration represents service
// configuration for kubectl exec ssh services (fka sockets).
type KubectlExecSshServiceConfiguration struct {
	KubectlExecTargetType string `json:"kubectl_exec_target_type"`

	BaseKubectlExecTargetConfiguration

	// mutually exclusive fields below
	StandardKubectlExecTargetConfiguration *StandardKubectlExecTargetConfiguration `json:"standard_kubectl_exec_target_configuration,omitempty"`
	AwsEksKubectlExecTargetConfiguration   *AwsEksKubectlExecTargetConfiguration   `json:"aws_eks_kubectl_exec_target_configuration,omitempty"`
}

// BaseKubectlExecTargetConfiguration represents base configuration for kubectl exec
// services (fka sockets), i.e. this configuration is common regardless of how the k8s
// cluster is hosted (aws, on prem, kind, etc...).
type BaseKubectlExecTargetConfiguration struct {
	// slice of allowed namespaces.
	NamespaceAllowlist []string `json:"namespace_allowlist,omitempty"`

	// map of namespace to selectors in that namespace.
	NamespaceSelectorsAllowlist map[string]map[string][]string `json:"namespace_selectors_allowlist,omitempty"`
}

// StandardKubectlExecTargetConfiguration represents service
// configuration for standard kubectl exec ssh services (fka sockets).
type StandardKubectlExecTargetConfiguration struct {
	MasterUrl      string `json:"master_url,omitempty"`
	KubeconfigPath string `json:"kubeconfig_path,omitempty"`
}

// AwsEksKubectlExecTargetConfiguration represents service
// configuration for aws eks kubectl exec ssh services (fka sockets).
type AwsEksKubectlExecTargetConfiguration struct {
	EksClusterName   string                 `json:"eks_cluster_name"`
	EksClusterRegion string                 `json:"eks_cluster_region"`
	AwsCredentials   *common.AwsCredentials `json:"aws_credentials,omitempty"`
}

// BuiltInSshServiceConfiguration represents the service configuration
// for the connector built-in ssh services (fka sockets).
type BuiltInSshServiceConfiguration struct {
	UsernameProvider string `json:"username_provider,omitempty"`
	Username         string `json:"username,omitempty"`
}

// UsernameAndPasswordAuthConfiguration represents authentication configuration
// for standard ssh services that require a username and password for authentication.
type UsernameAndPasswordAuthConfiguration struct {
	UsernameProvider string `json:"username_provider,omitempty"`
	Username         string `json:"username,omitempty"`
	Password         string `json:"password"`
}

// PrivateKeyAuthConfiguration represents authentication configuration
// for standard ssh services that require a private key for authentication.
type PrivateKeyAuthConfiguration struct {
	UsernameProvider string `json:"username_provider,omitempty"`
	Username         string `json:"username,omitempty"`
	PrivateKey       string `json:"private_key"`
}

// UsernameAndPasswordAuthConfiguration represents authentication configuration
// for standard ssh services that require a border0-signed certificate for authentication.
type Border0CertificateAuthConfiguration struct {
	UsernameProvider string `json:"username_provider,omitempty"`
	Username         string `json:"username,omitempty"`
}

// Validate validates the SshServiceConfiguration.
func (c *SshServiceConfiguration) Validate() error {
	switch c.SshServiceType {

	case SshServiceTypeAwsEc2InstanceConnect:
		if nilcheck.AnyNotNil(c.AwsSsmSshServiceConfiguration, c.StandardSshServiceConfiguration, c.BuiltInSshServiceConfiguration, c.DockerExecSshServiceConfiguration, c.KubectlExecSshServiceConfiguration) {
			return fmt.Errorf(
				"ssh service type \"%s\" can only have aws ec2 instance connect ssh service configuration defined",
				SshServiceTypeAwsEc2InstanceConnect)
		}
		if c.AwsEc2ICSshServiceConfiguration == nil {
			return fmt.Errorf(
				"ssh service configuration for ssh service type \"%s\" must have aws ec2 instance connect ssh service configuration defined",
				SshServiceTypeAwsEc2InstanceConnect,
			)
		}
		if err := c.AwsEc2ICSshServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid aws ec2 instance connect service configuration: %v", err)
		}
		return nil

	case SshServiceTypeAwsSsm:
		if nilcheck.AnyNotNil(c.AwsEc2ICSshServiceConfiguration, c.StandardSshServiceConfiguration, c.BuiltInSshServiceConfiguration, c.DockerExecSshServiceConfiguration, c.KubectlExecSshServiceConfiguration) {
			return fmt.Errorf(
				"ssh service type \"%s\" can only have aws ssm ssh service configuration defined",
				SshServiceTypeAwsSsm)
		}
		if c.AwsSsmSshServiceConfiguration == nil {
			return fmt.Errorf(
				"ssh service configuration for ssh service type \"%s\" must have aws ssm ssh service configuration defined",
				SshServiceTypeAwsSsm,
			)
		}
		if err := c.AwsSsmSshServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid aws ssm service configuration: %v", err)
		}
		return nil

	case SshServiceTypeConnectorBuiltIn:
		if nilcheck.AnyNotNil(c.AwsEc2ICSshServiceConfiguration, c.AwsSsmSshServiceConfiguration, c.StandardSshServiceConfiguration, c.DockerExecSshServiceConfiguration, c.KubectlExecSshServiceConfiguration) {
			return fmt.Errorf(
				"ssh service type \"%s\" can only have built in ssh service configuration defined",
				SshServiceTypeConnectorBuiltIn)
		}
		if c.BuiltInSshServiceConfiguration == nil {
			return fmt.Errorf(
				"ssh service configuration for ssh service type \"%s\" must have built in ssh service configuration defined",
				SshServiceTypeConnectorBuiltIn,
			)
		}
		if err := c.BuiltInSshServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid built in ssh service configuration: %v", err)
		}
		return nil

	case SshServiceTypeStandard:
		if nilcheck.AnyNotNil(c.AwsEc2ICSshServiceConfiguration, c.AwsSsmSshServiceConfiguration, c.BuiltInSshServiceConfiguration, c.DockerExecSshServiceConfiguration, c.KubectlExecSshServiceConfiguration) {
			return fmt.Errorf(
				"ssh service type \"%s\" can only have standard ssh service configuration defined",
				SshServiceTypeStandard)
		}
		if c.StandardSshServiceConfiguration == nil {
			return fmt.Errorf(
				"ssh service configuration for ssh service type \"%s\" must have standard ssh service configuration defined",
				SshServiceTypeStandard,
			)
		}
		if err := c.StandardSshServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid standard ssh service configuration: %v", err)
		}
		return nil

	case SshServiceTypeDockerExec:
		if nilcheck.AnyNotNil(c.AwsEc2ICSshServiceConfiguration, c.AwsSsmSshServiceConfiguration, c.BuiltInSshServiceConfiguration, c.KubectlExecSshServiceConfiguration, c.StandardSshServiceConfiguration) {
			return fmt.Errorf(
				"ssh service type \"%s\" can only have docker exec ssh service configuration defined",
				SshServiceTypeDockerExec)
		}
		// docker exec can be nil for now
		if c.DockerExecSshServiceConfiguration != nil {
			if err := c.DockerExecSshServiceConfiguration.Validate(); err != nil {
				return fmt.Errorf("invalid docker exec ssh service configuration: %v", err)
			}
		}
		return nil

	case SshServiceTypeKubectlExec:
		if nilcheck.AnyNotNil(c.AwsEc2ICSshServiceConfiguration, c.AwsSsmSshServiceConfiguration, c.BuiltInSshServiceConfiguration, c.DockerExecSshServiceConfiguration, c.StandardSshServiceConfiguration) {
			return fmt.Errorf(
				"ssh service type \"%s\" can only have kubectl exec ssh service configuration defined",
				SshServiceTypeKubectlExec)
		}
		if c.KubectlExecSshServiceConfiguration == nil {
			return fmt.Errorf(
				"ssh service configuration for ssh service type \"%s\" must have kubectl exec ssh service configuration defined",
				SshServiceTypeKubectlExec,
			)
		}
		if err := c.KubectlExecSshServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid kubectl exec ssh service configuration: %v", err)
		}
		return nil

	default:
		return fmt.Errorf("ssh service configuration has invalid ssh service type \"%s\"", c.SshServiceType)
	}
}

// Validate validates the AwsEc2ICSshServiceConfiguration.
func (c *AwsEc2ICSshServiceConfiguration) Validate() error {
	if err := c.HostnameAndPort.Validate(); err != nil {
		return err
	}
	if err := validateUsernameWithProvider(
		c.UsernameProvider,
		c.Username,
		set.New(UsernameProviderPromptClient),
	); err != nil {
		return err
	}
	if c.Ec2InstanceConnectEndpointId != "" {
		if !regexp.MustCompile(ec2InstanceConnectEndpointIdRegex).MatchString(c.Ec2InstanceConnectEndpointId) {
			return fmt.Errorf(
				"invalid ec2_instance_connect_endpoint_id: \"%s\" does not match regex \"%s\"",
				c.Ec2InstanceConnectEndpointId,
				ec2InstanceConnectEndpointIdRegex,
			)
		}
	}
	if c.Ec2InstanceId == "" {
		return fmt.Errorf("ec2_instance_id is a required field")
	}
	if c.Ec2InstanceRegion == "" {
		return fmt.Errorf("ec2_instance_region is a required field")
	}
	if err := common.ValidateAwsRegions(c.Ec2InstanceRegion); err != nil {
		return fmt.Errorf("invalid ec2_instance_region: %s", err)
	}
	if c.AwsCredentials != nil {
		if err := c.AwsCredentials.Validate(); err != nil {
			return fmt.Errorf("invalid aws_credentials: %v", err)
		}
	}
	return nil
}

// Validate validates the AwsSsmSshServiceConfiguration.
func (c *AwsSsmSshServiceConfiguration) Validate() error {
	switch c.SsmTargetType {

	case SsmTargetTypeEc2:
		if nilcheck.AnyNotNil(c.AwsSsmEcsTargetConfiguration) {
			return fmt.Errorf("ssm services with ssm target type \"%s\" can only have ec2 target configuration defined", SsmTargetTypeEc2)
		}
		if c.AwsSsmEc2TargetConfiguration == nil {
			return fmt.Errorf("ssm ec2 target configuration is required when ssm target type is \"%s\"", c.SsmTargetType)
		}
		if err := c.AwsSsmEc2TargetConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid aws ssm ec2 target configuration: %v", err)
		}
		return nil

	case SsmTargetTypeEcs:
		if nilcheck.AnyNotNil(c.AwsSsmEc2TargetConfiguration) {
			return fmt.Errorf("ssm services with ssm target type \"%s\" can only have ecs target configuration defined", SsmTargetTypeEcs)
		}
		if c.AwsSsmEcsTargetConfiguration == nil {
			return fmt.Errorf("ssm ecs target configuration is required when ssm target type is \"%s\"", c.SsmTargetType)
		}
		if err := c.AwsSsmEcsTargetConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid aws ssm ecs target configuration: %v", err)
		}
		return nil

	default:
		return fmt.Errorf("invalid ssm target type \"%s\"", c.SsmTargetType)
	}
}

// Validate validates the BuiltInSshServiceConfiguration.
func (c *BuiltInSshServiceConfiguration) Validate() error {
	return validateUsernameWithProvider(
		c.UsernameProvider,
		c.Username,
		set.New(
			UsernameProviderPromptClient,
			UsernameProviderUseConnectorUser,
		),
	)
}

func (c *StandardSshServiceConfiguration) Validate() error {
	if err := c.HostnameAndPort.Validate(); err != nil {
		return err
	}

	switch c.SshAuthenticationType {
	case StandardSshServiceAuthenticationTypeBorder0Certificate, "":
		if nilcheck.AnyNotNil(c.PrivateKeyAuthConfiguration, c.UsernameAndPasswordAuthConfiguration) {
			return fmt.Errorf(
				"ssh authentication type \"%s\" can only have border0 certificate auth configuration defined",
				StandardSshServiceAuthenticationTypeBorder0Certificate,
			)
		}
		if c.Border0CertificateAuthConfiguration == nil {
			return fmt.Errorf(
				"border0 certificate auth configuration is required when the ssh authentication type is \"%s\"",
				StandardSshServiceAuthenticationTypeBorder0Certificate,
			)
		}
		if err := c.Border0CertificateAuthConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid border0 certificate auth configuration: %v", err)
		}
		return nil

	case StandardSshServiceAuthenticationTypePrivateKey:
		if nilcheck.AnyNotNil(c.Border0CertificateAuthConfiguration, c.UsernameAndPasswordAuthConfiguration) {
			return fmt.Errorf(
				"ssh authentication type \"%s\" can only have private key auth configuration defined",
				StandardSshServiceAuthenticationTypePrivateKey,
			)
		}
		if c.PrivateKeyAuthConfiguration == nil {
			return fmt.Errorf(
				"private key auth configuration is required when the ssh authentication type is \"%s\"",
				StandardSshServiceAuthenticationTypePrivateKey,
			)
		}
		if err := c.PrivateKeyAuthConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid private key auth configuration: %v", err)
		}
		return nil

	case StandardSshServiceAuthenticationTypeUsernameAndPassword:
		if nilcheck.AnyNotNil(c.Border0CertificateAuthConfiguration, c.PrivateKeyAuthConfiguration) {
			return fmt.Errorf(
				"ssh authentication type \"%s\" can only have username and password auth configuration defined",
				StandardSshServiceAuthenticationTypeUsernameAndPassword,
			)
		}
		if c.UsernameAndPasswordAuthConfiguration == nil {
			return fmt.Errorf(
				"username and password auth configuration is required when the ssh authentication type is \"%s\"",
				StandardSshServiceAuthenticationTypeUsernameAndPassword,
			)
		}
		if err := c.UsernameAndPasswordAuthConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid username and password auth configuration: %v", err)
		}
		return nil

	default:
		return fmt.Errorf("invalid value for ssh_authentication_type: %s", c.SshAuthenticationType)
	}
}

// Validate validates a DockerExecSshServiceConfiguration.
func (c *DockerExecSshServiceConfiguration) Validate() error {
	regex := regexp.MustCompile(`^[a-zA-Z0-9*][a-zA-Z0-9*._\-]*$`)
	entries := set.New[string]()
	if len(c.ContainerNameAllowlist) > 0 {
		for i, name := range c.ContainerNameAllowlist {
			// reject empty string
			if name == "" {
				return fmt.Errorf("the container name allowlist entry in index %d is an empty string", i)
			}
			// make sure its valid
			if !regex.MatchString(name) {
				return fmt.Errorf("the container name allowlist entry in index %d (\"%s\") has invalid characters", i, name)
			}
			// make sure its not repeated
			if entries.Has(name) {
				return fmt.Errorf("the container name allowlist entry in index %d (\"%s\") is repeated", i, name)
			}
			entries.Add(name)
		}
	}
	return nil
}

// Validate validates a KubectlExecSshServiceConfiguration.
func (c *KubectlExecSshServiceConfiguration) Validate() error {
	// note: c.BaseKubectlExecTargetConfiguration is always valid

	switch c.KubectlExecTargetType {
	case KubectlExecTargetTypeStandard, "":
		if nilcheck.AnyNotNil(c.AwsEksKubectlExecTargetConfiguration) {
			return fmt.Errorf("kubectl exec ssh services with kubectl exec target type \"%s\" can only have standard kubectl exec target configuration defined", KubectlExecTargetTypeStandard)
		}
		// note: c.StandardKubectlExecTargetConfiguration can be nil
		if c.StandardKubectlExecTargetConfiguration != nil {
			if err := c.StandardKubectlExecTargetConfiguration.Validate(); err != nil {
				return fmt.Errorf("invalid standard kubectl exec target configuration: %v", err)
			}
		}
		return nil

	case KubectlExecTargetTypeAwsEks:
		if nilcheck.AnyNotNil(c.StandardKubectlExecTargetConfiguration) {
			return fmt.Errorf("kubectl exec ssh services with kubectl exec target type \"%s\" can only have aws eks kubectl exec target configuration defined", KubectlExecTargetTypeAwsEks)
		}
		if c.AwsEksKubectlExecTargetConfiguration == nil {
			return fmt.Errorf("aws eks kubectl exec target configuration is required when kubectl exec target type is \"%s\"", KubectlExecTargetTypeAwsEks)
		}
		if err := c.AwsEksKubectlExecTargetConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid aws eks kubectl exec target configuration: %v", err)
		}
		return nil

	default:
		return fmt.Errorf("invalid value for kubectl_exec_target_type: %s", c.KubectlExecTargetType)
	}
}

// Validate validates a StandardKubectlExecTargetConfiguration.
func (c *StandardKubectlExecTargetConfiguration) Validate() error {
	if c.MasterUrl != "" {
		if _, err := url.Parse(c.MasterUrl); err != nil {
			return fmt.Errorf("invalid value for master_url, invalid URL: %v", err)
		}
	}
	// note: can't really validate c.KubeconfigPath
	return nil
}

// Validate validates a AwsEksKubectlExecTargetConfiguration.
func (c *AwsEksKubectlExecTargetConfiguration) Validate() error {
	if c.EksClusterName == "" {
		return fmt.Errorf("eks_cluster_name is a required field")
	}
	if c.EksClusterRegion == "" {
		return fmt.Errorf("eks_cluster_region is a required field")
	}
	if err := common.ValidateAwsRegions(c.EksClusterRegion); err != nil {
		return fmt.Errorf("invalid eks_cluster_region: %s", err)
	}
	if c.AwsCredentials != nil {
		if err := c.AwsCredentials.Validate(); err != nil {
			return fmt.Errorf("invalid aws_credentials: %v", err)
		}
	}
	return nil
}

// Validate validates the AwsSsmEc2TargetConfiguration.
func (c *AwsSsmEc2TargetConfiguration) Validate() error {
	if c.Ec2InstanceId == "" {
		return fmt.Errorf("ec2_instance_id is a required field")
	}
	if c.Ec2InstanceRegion == "" {
		return fmt.Errorf("ec2_instance_region is a required field")
	}
	if err := common.ValidateAwsRegions(c.Ec2InstanceRegion); err != nil {
		return fmt.Errorf("invalid ec2_instance_region: %s", err)
	}
	if c.AwsCredentials != nil {
		if err := c.AwsCredentials.Validate(); err != nil {
			return fmt.Errorf("invalid aws_credentials: %v", err)
		}
	}
	return nil
}

// Validate validates the AwsSsmEcsTargetConfiguration.
func (c *AwsSsmEcsTargetConfiguration) Validate() error {
	if c.EcsClusterName == "" {
		return fmt.Errorf("ecs_cluster_name is a required field")
	}
	if c.EcsClusterRegion == "" {
		return fmt.Errorf("ecs_cluster_region is a required field")
	}
	if err := common.ValidateAwsRegions(c.EcsClusterRegion); err != nil {
		return fmt.Errorf("invalid ecs_cluster_region: %s", err)
	}
	if c.EcsServiceName == "" {
		return fmt.Errorf("ecs_service_name is a required field")
	}
	if c.AwsCredentials != nil {
		if err := c.AwsCredentials.Validate(); err != nil {
			return fmt.Errorf("invalid aws_credentials: %v", err)
		}
	}
	return nil
}

// Validate validates the Border0CertificateAuthConfiguration.
func (c *Border0CertificateAuthConfiguration) Validate() error {
	return validateUsernameWithProvider(
		c.UsernameProvider,
		c.Username,
		set.New(UsernameProviderPromptClient),
	)
}

// Validate validates the PrivateKeyAuthConfiguration.
func (c *PrivateKeyAuthConfiguration) Validate() error {
	if err := validateUsernameWithProvider(
		c.UsernameProvider,
		c.Username,
		set.New[string](),
	); err != nil {
		return err
	}
	if c.PrivateKey == "" {
		return fmt.Errorf("private_key is a required field")
	}
	if !strings.HasPrefix(c.PrivateKey, "from:") {
		_, err := ssh.ParseRawPrivateKey([]byte(c.PrivateKey))
		if err != nil {
			return fmt.Errorf("private_key is not a valid PEM or DER encoded private key")
		}
	}

	return nil
}

// Validate validates the Border0CertificateAuthConfiguration.
func (c *UsernameAndPasswordAuthConfiguration) Validate() error {
	if err := validateUsernameWithProvider(
		c.UsernameProvider,
		c.Username,
		set.New[string](),
	); err != nil {
		return err
	}
	if c.Password == "" {
		return fmt.Errorf("password is a required field")
	}
	return nil
}
