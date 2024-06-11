package service

import (
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"

	"github.com/borderzero/border0-go/lib/types/empty"
)

// KubernetesServiceConfiguration represents service
// configuration for kubernetes services (fka sockets).
type KubernetesServiceConfiguration struct {
	Server                   string `json:"server,omitempty"`
	CertificateAuthorityData string `json:"certificate_authority_data,omitempty"`
	ClientCertificateData    string `json:"client_certificate_data,omitempty"`
	ClientKeyData            string `json:"client_key_data,omitempty"`
}

// Validate validates the KubernetesServiceConfiguration.
func (c *KubernetesServiceConfiguration) Validate() error {
	if _, err := url.Parse(c.Server); err != nil {
		return fmt.Errorf("invalid value for server, invalid URL: %v", err)
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
	}
	allOrNoneFields := []string{
		c.CertificateAuthorityData,
		c.ClientCertificateData,
		c.ClientKeyData,
	}
	if n := empty.Count(allOrNoneFields...); n != 0 && n != len(allOrNoneFields) {
		return fmt.Errorf("either all or none of certificate_authority_data, client_certificate_data, and client_key_data must be provided, got %d/%d", len(allOrNoneFields)-n, len(allOrNoneFields))
	}
	return nil
}
