package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateKubernetesServiceConfiguration(t *testing.T) {
	t.Parallel()

	validServer := "https://127.0.0.1:55606"
	validCertificateAuthorityData := "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURLVENDQWhHZ0F3SUJBZ0lJZTcwRTE3cUprajh3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TkRBMk1EVXhPREExTVRoYUZ3MHlOVEEyTURVeE9ERXdNakJhTUR3eApIekFkQmdOVkJBb1RGbXQxWW1WaFpHMDZZMngxYzNSbGNpMWhaRzFwYm5NeEdUQVhCZ05WQkFNVEVHdDFZbVZ5CmJtVjBaWE10WVdSdGFXNHdnZ0VpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFEcElmOCsKQ2RuZlRSeFU1a1FtUEZvVUhZNkNrcUR3bFZPK1BMdWJXSnFuS3NHTW1Cdmp4R0pjSEtRaFlYVWNPKzRzZi9ocQpCcHFWZ3dhcHRRcHE3eWxhM0gwTzhwWGJXdmpsdm1xcUtmV1hsUlZzd1ZwUkhSR1QzWGFuUVVoRDVYYVF0UUNlCktOL2FheERMUXJzMW8xOXFoNTRSOERHd2NxaktqdmNnaTMzQjcwbVBoZHZxU1JlYzBGN3FmSjM2bjFUQUFDYmgKdkNJWDZtQ0Noa1RLeGNSRm9Mb3grRU1iZXZVdkQ1TUlNV0tQd3o3azAxTXhmNkNFVEdLK3YvZGs3cWJHOU1hcQoxTzdZS0JxaHhNVXF6ek03ZHBQVTVnZkFib0Nlek5OTHJmaGxFT2o2ZXZrbTFpQ1VrNHhVWEkySkxXUXpFcHFRCmtwRVZVRUtUNTI2V29QQy9BZ01CQUFHalZqQlVNQTRHQTFVZER3RUIvd1FFQXdJRm9EQVRCZ05WSFNVRUREQUsKQmdnckJnRUZCUWNEQWpBTUJnTlZIUk1CQWY4RUFqQUFNQjhHQTFVZEl3UVlNQmFBRkk5Y1JRNTVqQ2dDRDVTRQpGVDNRQkNoYStpY25NQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUFKamo3OTNGV1NFMm9hTHAweUtENFZnQXJLCk8rSllQcEpqZENwbEhsNk1GREpSWEt3RENQeTZXbDZhUkdWTUJmQUp5cXoyU1djNWdBeXZFTGRXTkVVaHhWaU0KQk9YRDZIMXRVWVg4ZE1DeWw2LzJ0NDNvRVpxdnVlR1ZXMjR6WllUTGpEOWwxQnBGRElYK2FjWHFsZlR5dWFoZApLbXBCN1ErQ0pSV3ZwTHJsNW4zM1RDSmhRbldlRDdXYnR3bjBtS2JGaDJjekJEQTVJUUM4djA5a3h5ZWhQOHlSCnYweWZSOFB2MDZXZzlxQlJxRGhZWDhZY2VOaEQ4ejd5S1NsQjdVWlZnVm9TcGEyL2Fra0pNaUNrNnA4dlBOZGEKSjZkcE1kREptMWxkcUtYTXlTWlhjOFVOOVJldzZhN2x0Unk5SDNvaHZEK1VwcWtYbnlSbTlOdUtnRVVXCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
	validClientCertificateData := "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBNlNIL1BnblozMDBjVk9aRUpqeGFGQjJPZ3BLZzhKVlR2ank3bTFpYXB5ckJqSmdiCjQ4UmlYQnlrSVdGMUhEdnVMSC80YWdhYWxZTUdxYlVLYXU4cFd0eDlEdktWMjFyNDViNXFxaW4xbDVVVmJNRmEKVVIwUms5MTJwMEZJUStWMmtMVUFuaWpmMm1zUXkwSzdOYU5mYW9lZUVmQXhzSEtveW83M0lJdDl3ZTlKajRYYgo2a2tYbk5CZTZueWQrcDlVd0FBbTRid2lGK3BnZ29aRXlzWEVSYUM2TWZoREczcjFMdytUQ0RGaWo4TSs1Tk5UCk1YK2doRXhpdnIvM1pPNm14dlRHcXRUdTJDZ2FvY1RGS3M4ek8zYVQxT1lId0c2QW5zelRTNjM0WlJEbytucjUKSnRZZ2xKT01WRnlOaVMxa014S2FrSktSRlZCQ2srZHVscUR3dndJREFRQUJBb0lCQUVHV3hQNnBPdGl1RXhoMwpxQml4eGczK3VzZFZoVjNxUlNNeVp2MnZrZU1TRGw0RnRXTUJMME1HODg1SVl2Qk1sQmVFOFZpeS91NnRiRVIyCm9MV3hENHQybU9DSnMwdUJicUVFVDNJTVdBSUxlWG1kZEZGNXdpWTltbEkzOWFMVWZlVDhaMjRYeGRMbmFGUW4KTkJCRkpHVU9QcFNCWlJsbFNNT2tNVGJBeCtLMWJBTHlmYUdXMmlVWlZ5a0hYT3E0L3p1SVlOeHFDL1BtcDYrTgpZZzJWZWkrOFN6OGxaUFZtUHNqd1ZjV24xTENYd0JwaE1wVmU2MEFTN1EzZFNkY3pmQnVlZ3NBRjliRWFGWGdXCm4wK1htSTNNTW1IdVdSRzk1TzlRdFNWM0hKWTRNMVJsSE5PNnVCazAwMDdTMUx4YVZVMWE3VU41eFpkYlk3SlYKRGtkanVzRUNnWUVBNm91UERCVlhwTDNDcTQ3T3FFci9ZWWd5bzVxblFweUMyd3dvUTIyY0pGcjlheDdGK2U5dQpCL1J2eXZRdER0SHd6c2hlcnU5YkViK3pPYmYyQ3d4czhDWjhGODU1aEZoMHBTUUFDQ0FxNVZudFZnSWVkTDFlClBKb3pLekkybzdEeSs4d0krVmRpQjFFSWZKNm5wMzdBN1BJNHhnNThNaW5GVXZZajlKbk5EcThDZ1lFQS9uVmQKV3BBMWRVS1ZkVmhzWHVhTmRPc2NaZmg3d0VWTmJLbDVnSG5NZHJmN0o5RlpESzVFNStpaW9jSHBCcEQ0S240OQpNNjBQTWVWQWVxQzc0RmZFTEc4c1hwTHRPODRLTER3bVUzMVhoWW00anBLMCs2dk5BQnRLMEpWNlRTam11bjdHCkxWYmpVSnVqL1pUbzZtSk1Ba3UwamsyZUxFRVZwMnZoWnhsRFF2RUNnWUVBcEdBT2JwOStmdnhtdENrdmVBNm4Ka0VrTmhFOWNyWGhXMVFGZUxuTmhITEdRTFVTeDV3b0FDUjhzUWdhQ2xZSll5L0x0T0ppVE1Jc3pKbDVnMXpscwpGNXBCR2NZZ3I3bjJkYzRYSDVxN3RXVWl5a1pONWU3WlhvdUxGTmpxSmlwYkFGRHNjU0xtbHM0WnJvemFYcTZjCnN2TEhDemNYbjJoYTNGdVlzMUlldmcwQ2dZRUFsTGZGZ3V1cndzZlNrRTNJc01Kd3hIZXhGa3ZmRCtXeE5hUHIKVmxKWlZMMTh0YTNlTE5JaGhhVHFnRDNUMHJtaHdUd2N5dm9IV25NUmQxby8xQU9YclJ5Z3d2bCtXNmkxTzF0QQpUeVYvcEpWMFh5ZVJUVklBczJKQUhKNXdaMjVUUTRaYWp4OHNKZTJCTC9EN0hCbXRNTjVNUGF0Um1hM0VXU0J0CkVaT3JReEVDZ1lBL0xEOVRxY1FFUTFtOG5DV2xvNnNXV05DbEZ5MWhIcEpwZ1E3ZHRKZjI1akdmQW9QREpnRG8KcnpnZlFMNnNMNkJ2a29tc2lwNC9BMXp5NERwOGJzQVI5Q1daT0NoV1dYdHpCM2NtbFRUbys1bGFSVDRUM2dLRQoyS09hN1JzU2FYVzFKMWpralQzL1I2YnA4cXNkbVVVSVZTcDAzdnJhY0laS1kzRnc3UkU3dXc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
	validClientKeyData := "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURCVENDQWUyZ0F3SUJBZ0lJTmZHNUlHTHVFYm93RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TkRBMk1EVXhPREExTVRoYUZ3MHpOREEyTURNeE9ERXdNVGhhTUJVeApFekFSQmdOVkJBTVRDbXQxWW1WeWJtVjBaWE13Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLCkFvSUJBUUM4YU1kZ3gxN0FrZVhPTGwxNUQ2c0FPNndkWnhmVzNFa1MzajMxWXI5UzdUMTVsRVczSktmbzZPMzYKTlhTVVREZFVmRXlhQ0tNQTR0bjlaRkZ1aG9jNUljQXBlc09ZcExKS3NjUk01WCtFSmVWdWg2Rzd4TThyNVh1ZAp5VExEOHFTS3Z2UjZQcmxhQnMzSWUvTlV5aVBvUVRpMWJLcTFzMG5uZ1pEWHlnbTJ3bFFxemhlS3hJWmdqTUtyCjI4Y2hsWUtucGpNYXllREF2OGVtOG41dkJnWXl6OHBmYXBaWnQ1RTJ5a2hwVnBpYzVOTm1jbEdGVklxbitCYjIKSXZ1bWZIelRJc1h2cVlQOFVBMU5EMFl6VzJDcWlwNEFJWEd2V040OC8xeVgwQjRLNHlCUkNJbzdzekdoSzJQKwpUeHhLc2NwMlR3Vlp4VUlZKzJpR0t6L1dmZkQzQWdNQkFBR2pXVEJYTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQCkJnTlZIUk1CQWY4RUJUQURBUUgvTUIwR0ExVWREZ1FXQkJTUFhFVU9lWXdvQWcrVWhCVTkwQVFvV3Zvbkp6QVYKQmdOVkhSRUVEakFNZ2dwcmRXSmxjbTVsZEdWek1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQ1ZwcUJTUUxybwpBc3pncUtmc0RMZVZhSHBROXhiQndaWXB3Tm1tZDI2bnB1SXFEVnZzU29LbW5HVXFqV2phT24rUkMrMkdUNndrClcyU2RPWU8rQkd4TGZ6UE5zKzhlOE1KVGtUMWRTSytSaURuOWpHaDJnZVVBbGh4NlFEZjY0aGFYWDZBT3V0aTUKdFl0VE42ZkdFdERDV3lNTnAyY1o0ZFJEbzZiMHFMUkE0QlZzRzU0RFVMNWFkaVMvM0NCZW4yZTdZeHV4bWlBago1elF5V01QQU1jQTV2YzYvQUs2ZkxtYmFpanNuRytoVldNRGZhb05wbzI4MkdZYXBKVWJscTk5aVE5TDVqUlVqCnd2VUdHYnlSZmdWdmdrZEliaGFUd3B4OFQrRnZXNVEvWnlmaStFZW5uYXBXRUQwOU1sUE5uZDR4UFRJbCtPSTYKdTdZSmNSTWF3UVV4Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
	invalidDataNotB64 := "!"                // not a valid base64 encoding character
	invalidDataNotPEM := "aGVsbG8gd29ybGQK" // just base64 encoded "hello world"

	tests := []struct {
		name          string
		configuration *KubernetesServiceConfiguration
		expectedError error
	}{
		{
			name: "Happy case for kubernetes service with no cert data",
			configuration: &KubernetesServiceConfiguration{
				Server: validServer,
			},
			expectedError: nil,
		},
		{
			name: "Happy case for kubernetes service with cert data",
			configuration: &KubernetesServiceConfiguration{
				Server:                   validServer,
				CertificateAuthorityData: validCertificateAuthorityData,
				ClientCertificateData:    validClientCertificateData,
				ClientKeyData:            validClientKeyData,
			},
			expectedError: nil,
		},
		{
			name: "Failure case for kubernetes service with partial cert data - missing ca cert",
			configuration: &KubernetesServiceConfiguration{
				Server:                validServer,
				ClientCertificateData: validClientCertificateData,
				ClientKeyData:         validClientKeyData,
			},
			expectedError: errors.New("either all or none of certificate_authority_data, client_certificate_data, and client_key_data must be provided, got 2/3"),
		},
		{
			name: "Failure case for kubernetes service with partial cert data - missing client cert",
			configuration: &KubernetesServiceConfiguration{
				Server:                   validServer,
				CertificateAuthorityData: validCertificateAuthorityData,
				ClientKeyData:            validClientKeyData,
			},
			expectedError: errors.New("either all or none of certificate_authority_data, client_certificate_data, and client_key_data must be provided, got 2/3"),
		},
		{
			name: "Failure case for kubernetes service with partial cert data - missing client cert key",
			configuration: &KubernetesServiceConfiguration{
				Server:                   validServer,
				CertificateAuthorityData: validCertificateAuthorityData,
				ClientCertificateData:    validClientCertificateData,
			},
			expectedError: errors.New("either all or none of certificate_authority_data, client_certificate_data, and client_key_data must be provided, got 2/3"),
		},
		{
			name: "Failure case for kubernetes service with bad cert data - bad ca cert base64",
			configuration: &KubernetesServiceConfiguration{
				Server:                   validServer,
				CertificateAuthorityData: invalidDataNotB64,
				ClientCertificateData:    validClientCertificateData,
				ClientKeyData:            validClientKeyData,
			},
			expectedError: errors.New("failed to base64-decode certificate_authority_data: illegal base64 data at input byte 0"),
		},
		{
			name: "Failure case for kubernetes service with bad cert data - bad client cert base64",
			configuration: &KubernetesServiceConfiguration{
				Server:                   validServer,
				CertificateAuthorityData: validCertificateAuthorityData,
				ClientCertificateData:    invalidDataNotB64,
				ClientKeyData:            validClientKeyData,
			},
			expectedError: errors.New("failed to base64-decode client_certificate_data: illegal base64 data at input byte 0"),
		},
		{
			name: "Failure case for kubernetes service with bad cert data - bad client cert key base64",
			configuration: &KubernetesServiceConfiguration{
				Server:                   validServer,
				CertificateAuthorityData: validCertificateAuthorityData,
				ClientCertificateData:    validClientCertificateData,
				ClientKeyData:            invalidDataNotB64,
			},
			expectedError: errors.New("failed to base64-decode client_key_data: illegal base64 data at input byte 0"),
		},
		{
			name: "Failure case for kubernetes service with bad cert data - bad ca cert PEM",
			configuration: &KubernetesServiceConfiguration{
				Server:                   validServer,
				CertificateAuthorityData: invalidDataNotPEM,
				ClientCertificateData:    validClientCertificateData,
				ClientKeyData:            validClientKeyData,
			},
			expectedError: errors.New("failed to PEM-decode certificate_authority_data: not valid PEM"),
		},
		{
			name: "Failure case for kubernetes service with bad cert data - bad client cert PEM",
			configuration: &KubernetesServiceConfiguration{
				Server:                   validServer,
				CertificateAuthorityData: validCertificateAuthorityData,
				ClientCertificateData:    invalidDataNotPEM,
				ClientKeyData:            validClientKeyData,
			},
			expectedError: errors.New("failed to PEM-decode client_certificate_data: not valid PEM"),
		},
		{
			name: "Failure case for kubernetes service with bad cert data - bad ca cert key PEM",
			configuration: &KubernetesServiceConfiguration{
				Server:                   validServer,
				CertificateAuthorityData: validCertificateAuthorityData,
				ClientCertificateData:    validClientCertificateData,
				ClientKeyData:            invalidDataNotPEM,
			},
			expectedError: errors.New("failed to PEM-decode client_key_data: not valid PEM"),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.configuration.Validate())
		})
	}
}
