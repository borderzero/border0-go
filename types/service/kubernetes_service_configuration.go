package service

import (
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"

	"github.com/borderzero/border0-go/lib/types/nilcheck"
	"github.com/borderzero/border0-go/types/common"
)

// Kubernetes service types supported by Border0. Choose `standard` for
// self-managed kubernetes clusters  or `aws_eks` for AWS EKS clusters.
const (
	KubernetesServiceTypeStandard = "standard"
	KubernetesServiceTypeAwsEks   = "aws_eks"
)

// KubernetesServiceConfiguration represents service
// configuration for kubernetes services (fka sockets).
type KubernetesServiceConfiguration struct {
	KubernetesServiceType string `json:"kubernetes_service_type"`

	// mutually exclusive fields below
	StandardKubernetesServiceConfiguration *StandardKubernetesServiceConfiguration `json:"standard_kubernetes_service_configuration,omitempty"`
	AwsEksKubernetesServiceConfiguration   *AwsEksKubernetesServiceConfiguration   `json:"aws_eks_kubernetes_service_configuration,omitempty"`
}

// StandardKubernetesServiceConfiguration represents service
// configuration for standard kubernetes services (fka sockets).
type StandardKubernetesServiceConfiguration struct {
	// for the connector to load config from the filesystem
	KubeconfigPath string `json:"kubeconfig_path,omitempty"`
	Context        string `json:"context,omitempty"`

	// for the connector to communicate to the kubernetes api server
	Server                   string `json:"server,omitempty"`
	CertificateAuthority     string `json:"certificate_authority,omitempty"`
	CertificateAuthorityData string `json:"certificate_authority_data,omitempty"`

	// for the kubernetes api server to authenticate the connector with client certificates
	ClientCertificate     string `json:"client_certificate,omitempty"`
	ClientCertificateData string `json:"client_certificate_data,omitempty"`
	ClientKey             string `json:"client_key,omitempty"`
	ClientKeyData         string `json:"client_key_data,omitempty"`

	// for the kubernetes api server to authenticate the connector with a token
	Token     string `json:"token,omitempty"`
	TokenFile string `json:"token_file,omitempty"`
}

// AwsEksKubernetesServiceConfiguration represents service
// configuration for aws eks kubernetes services (fka sockets).
type AwsEksKubernetesServiceConfiguration struct {
	EksClusterName   string                 `json:"eks_cluster_name"`
	EksClusterRegion string                 `json:"eks_cluster_region"`
	AwsCredentials   *common.AwsCredentials `json:"aws_credentials,omitempty"`
}

// Validate validates the KubernetesServiceConfiguration.
func (c *KubernetesServiceConfiguration) Validate() error {
	switch c.KubernetesServiceType {

	case KubernetesServiceTypeStandard:
		if nilcheck.AnyNotNil(c.AwsEksKubernetesServiceConfiguration) {
			return fmt.Errorf(
				"kubernetes service type \"%s\" can only have standard kubernetes service configuration defined",
				KubernetesServiceTypeStandard,
			)
		}
		if c.StandardKubernetesServiceConfiguration != nil {
			if err := c.StandardKubernetesServiceConfiguration.Validate(); err != nil {
				return fmt.Errorf("invalid standard kubernetes service configuration: %v", err)
			}
		}
		return nil

	case KubernetesServiceTypeAwsEks:
		if nilcheck.AnyNotNil(c.StandardKubernetesServiceConfiguration) {
			return fmt.Errorf(
				"kubernetes service type \"%s\" can only have aws eks kubernetes service configuration defined",
				KubernetesServiceTypeAwsEks,
			)
		}
		if c.AwsEksKubernetesServiceConfiguration == nil {
			return fmt.Errorf(
				"kubernetes service configuration for kubernetes service type \"%s\" must have aws eks kubernetes service configuration defined",
				KubernetesServiceTypeAwsEks,
			)
		}
		if err := c.AwsEksKubernetesServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid aws eks kubernetes service configuration: %v", err)
		}
		return nil

	default:
		return fmt.Errorf("kubernetes service configuration has invalid kubernetes service type \"%s\"", c.KubernetesServiceType)
	}
}

// Validate validates the StandardKubernetesServiceConfiguration.
func (c *StandardKubernetesServiceConfiguration) Validate() error {
	if c.Server != "" {
		if _, err := url.Parse(c.Server); err != nil {
			return fmt.Errorf("invalid value for server, invalid URL: %v", err)
		}
	}
	if c.CertificateAuthorityData != "" {
		pemBytes, err := base64.StdEncoding.DecodeString(c.CertificateAuthorityData)
		if err != nil {
			return fmt.Errorf("failed to base64-decode certificate_authority_data: %v", err)
		}
		pemBlock, _ := pem.Decode(pemBytes)
		if pemBlock == nil {
			return errors.New("failed to PEM-decode certificate_authority_data: not valid PEM")
		}
	}
	if c.ClientCertificateData != "" {
		pemBytes, err := base64.StdEncoding.DecodeString(c.ClientCertificateData)
		if err != nil {
			return fmt.Errorf("failed to base64-decode client_certificate_data: %v", err)
		}
		pemBlock, _ := pem.Decode(pemBytes)
		if pemBlock == nil {
			return errors.New("failed to PEM-decode client_certificate_data: not valid PEM")
		}
		if c.ClientKey == "" && c.ClientKeyData == "" {
			return errors.New("client_certificate_data was provided but both client_key and client_key_data were empty")
		}
	}
	if c.ClientCertificate != "" {
		if c.ClientKey == "" && c.ClientKeyData == "" {
			return errors.New("client_certificate was provided but both client_key and client_key_data were empty")
		}
	}
	if c.ClientKeyData != "" {
		pemBytes, err := base64.StdEncoding.DecodeString(c.ClientKeyData)
		if err != nil {
			return fmt.Errorf("failed to base64-decode client_key_data: %v", err)
		}
		pemBlock, _ := pem.Decode(pemBytes)
		if pemBlock == nil {
			return errors.New("failed to PEM-decode client_key_data: not valid PEM")
		}
		if c.ClientCertificate == "" && c.ClientCertificateData == "" {
			return errors.New("client_key_data was provided but both client_certificate and client_certificate_data were empty")
		}
	}
	if c.ClientKey != "" {
		if c.ClientCertificate == "" && c.ClientCertificateData == "" {
			return errors.New("client_key was provided but both client_certificate and client_certificate_data were empty")
		}
	}
	return nil
}

// Validate validates the AwsEksKubernetesServiceConfiguration.
func (c *AwsEksKubernetesServiceConfiguration) Validate() error {
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
