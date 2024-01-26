package service

import (
	"errors"
	"testing"

	"github.com/borderzero/border0-go/lib/types/pointer"
	"github.com/borderzero/border0-go/types/common"
	"github.com/stretchr/testify/assert"
)

func Test_DatabaseServiceConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testStandardConfig := &StandardDatabaseServiceConfiguration{
		DatabaseProtocol: DatabaseProtocolMySql,
		HostnameAndPort: HostnameAndPort{
			Hostname: "hostname",
			Port:     3306,
		},
		AuthenticationType: DatabaseAuthenticationTypeTls,
		TlsAuth: &DatabaseTlsAuthConfiguration{
			UsernameAndPassword: UsernameAndPassword{
				Username: "username",
				Password: "password",
			},
			TlsConfig: TlsConfig{
				Certificate: "certificate",
				Key:         "private-key",
			},
		},
	}
	testAwsRdsConfig := &AwsRdsDatabaseServiceConfiguration{
		DatabaseProtocol: DatabaseProtocolMySql,
		HostnameAndPort: HostnameAndPort{
			Hostname: "hostname",
			Port:     3306,
		},
		AuthenticationType: DatabaseAuthenticationTypeIam,
		IamAuth: &AwsRdsIamAuthConfiguration{
			Username:          "username",
			RdsInstanceRegion: "us-east-1",
			AwsCredentials: &common.AwsCredentials{
				AwsAccessKeyId:     pointer.To("AKIA000FAKE00KEY00ID"),
				AwsSecretAccessKey: pointer.To("Secret+Access+Key/0000000000000000000000"),
			},
		},
	}
	testGcpCloudSqlConfig := &GcpCloudSqlDatabaseServiceConfiguration{
		CloudSqlConnectorEnabled: false,
		Standard: &GcpCloudSqlStandardConfiguration{
			DatabaseProtocol: DatabaseProtocolMySql,
			HostnameAndPort: HostnameAndPort{
				Hostname: "hostname",
				Port:     3306,
			},
			AuthenticationType: DatabaseAuthenticationTypeTls,
			TlsAuth: &DatabaseTlsAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: "username",
					Password: "password",
				},
				TlsConfig: TlsConfig{
					Certificate: "certificate",
					Key:         "private-key",
				},
			},
		},
	}

	tests := []struct {
		name  string
		given DatabaseServiceConfiguration
		want  error
	}{
		{
			name:  "database service type is missing",
			given: DatabaseServiceConfiguration{
				// database service type is missing
			},
			want: errors.New("database service type is required"),
		},
		{
			name: "invalid database service type",
			given: DatabaseServiceConfiguration{
				DatabaseServiceType: "invalid",
			},
			want: errors.New("invalid database service type: invalid"),
		},
		{
			name: "when standard type picked, other database service configs should be nil",
			given: DatabaseServiceConfiguration{
				DatabaseServiceType: DatabaseServiceTypeStandard,
				Standard:            testStandardConfig,
				AwsRds:              testAwsRdsConfig,
				GcpCloudSql:         testGcpCloudSqlConfig,
			},
			want: errors.New("database service type is standard, but AWS RDS, Google Cloud SQL or Azure SQL configuration is provided"),
		},
		{
			name: "standard type is picked, but standard config is missing",
			given: DatabaseServiceConfiguration{
				DatabaseServiceType: DatabaseServiceTypeStandard,
				// standard config is missing
			},
			want: errors.New("standard database service configuration is required"),
		},
		{
			name: "when aws rds type picked, other database service configs should be nil",
			given: DatabaseServiceConfiguration{
				DatabaseServiceType: DatabaseServiceTypeAwsRds,
				Standard:            testStandardConfig,
				AwsRds:              testAwsRdsConfig,
				GcpCloudSql:         testGcpCloudSqlConfig,
			},
			want: errors.New("database service type is aws_rds, but standard, Google Cloud SQL or Azure SQL configuration is provided"),
		},
		{
			name: "aws rds type is picked, but aws rds config is missing",
			given: DatabaseServiceConfiguration{
				DatabaseServiceType: DatabaseServiceTypeAwsRds,
				// aws rds config is missing
			},
			want: errors.New("AWS RDS database service configuration is required"),
		},
		{
			name: "when google cloud sql type picked, other database service configs should be nil",
			given: DatabaseServiceConfiguration{
				DatabaseServiceType: DatabaseServiceTypeGcpCloudSql,
				Standard:            testStandardConfig,
				AwsRds:              testAwsRdsConfig,
				GcpCloudSql:         testGcpCloudSqlConfig,
			},
			want: errors.New("database service type is gcp_cloudsql, but standard, AWS RDS or Azure SQL configuration is provided"),
		},
		{
			name: "google cloud sql type is picked, but google cloud sql config is missing",
			given: DatabaseServiceConfiguration{
				DatabaseServiceType: DatabaseServiceTypeGcpCloudSql,
				// google cloud sql config is missing
			},
			want: errors.New("Google Cloud SQL database service configuration is required"),
		},
		{
			name: "happy path - standard config",
			given: DatabaseServiceConfiguration{
				DatabaseServiceType: DatabaseServiceTypeStandard,
				Standard:            testStandardConfig,
			},
			want: nil,
		},
		{
			name: "happy path - aws rds config",
			given: DatabaseServiceConfiguration{
				DatabaseServiceType: DatabaseServiceTypeAwsRds,
				AwsRds:              testAwsRdsConfig,
			},
			want: nil,
		},
		{
			name: "happy path - google cloud sql config",
			given: DatabaseServiceConfiguration{
				DatabaseServiceType: DatabaseServiceTypeGcpCloudSql,
				GcpCloudSql:         testGcpCloudSqlConfig,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_StandardDatabaseServiceConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testDatbaseProtocol := DatabaseProtocolMySql
	testHostnameAndPort := HostnameAndPort{
		Hostname: "hostname",
		Port:     3306,
	}
	testPasswordAuth := &DatabaseUsernameAndPasswordAuthConfiguration{
		UsernameAndPassword: UsernameAndPassword{
			Username: "username",
			Password: "password",
		},
	}
	testTlsAuth := &DatabaseTlsAuthConfiguration{
		UsernameAndPassword: UsernameAndPassword{
			Username: "username",
			Password: "password",
		},
		TlsConfig: TlsConfig{
			Certificate: "certificate",
			Key:         "private-key",
		},
	}

	tests := []struct {
		name  string
		given StandardDatabaseServiceConfiguration
		want  error
	}{
		{
			name:  "database protocol is missing",
			given: StandardDatabaseServiceConfiguration{
				// database protocol is missing
			},
			want: errors.New("database protocol is required"),
		},
		{
			name: "hostname is missing",
			given: StandardDatabaseServiceConfiguration{
				DatabaseProtocol: testDatbaseProtocol,
				HostnameAndPort: HostnameAndPort{
					// hostname is missing
					Port: 3306,
				},
				AuthenticationType:      DatabaseAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuth: testPasswordAuth,
			},
			want: errors.New("hostname is a required field"),
		},
		{
			name: "invalid database authentication type",
			given: StandardDatabaseServiceConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: "invalid",
			},
			want: errors.New("invalid database authentication type: invalid"),
		},
		{
			name: "when username and password auth is picked, tls auth config should be nil",
			given: StandardDatabaseServiceConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				HostnameAndPort:         testHostnameAndPort,
				AuthenticationType:      DatabaseAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuth: testPasswordAuth,
				TlsAuth:                 testTlsAuth,
			},
			want: errors.New("authentication type is username_and_password, but tls_auth, kerberos or sql_authentication configuration is provided"),
		},
		{
			name: "username and password auth is picked, but username and password auth config is missing",
			given: StandardDatabaseServiceConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: DatabaseAuthenticationTypeUsernameAndPassword,
				// username and password auth config is missing
			},
			want: errors.New("username and password auth configuration is required"),
		},
		{
			name: "when tls auth is picked, username and password auth config should be nil",
			given: StandardDatabaseServiceConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				HostnameAndPort:         testHostnameAndPort,
				AuthenticationType:      DatabaseAuthenticationTypeTls,
				UsernameAndPasswordAuth: testPasswordAuth,
				TlsAuth:                 testTlsAuth,
			},
			want: errors.New("authentication type is tls, but username_and_password, kerberos or sql_authentication configuration is provided"),
		},
		{
			name: "tls auth is picked, but tls auth config is missing",
			given: StandardDatabaseServiceConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: DatabaseAuthenticationTypeTls,
				// tls auth config is missing
			},
			want: errors.New("TLS auth configuration is required"),
		},
		{
			name: "happy path - username and password auth",
			given: StandardDatabaseServiceConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				HostnameAndPort:         testHostnameAndPort,
				AuthenticationType:      DatabaseAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuth: testPasswordAuth,
			},
			want: nil,
		},
		{
			name: "happy path - tls auth",
			given: StandardDatabaseServiceConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: DatabaseAuthenticationTypeTls,
				TlsAuth:            testTlsAuth,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_AwsRdsDatabaseServiceConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testDatbaseProtocol := DatabaseProtocolMySql
	testHostnameAndPort := HostnameAndPort{
		Hostname: "hostname",
		Port:     3306,
	}
	testPasswordAuth := &AwsRdsUsernameAndPasswordAuthConfiguration{
		UsernameAndPassword: UsernameAndPassword{
			Username: "username",
			Password: "password",
		},
	}
	testIamAuth := &AwsRdsIamAuthConfiguration{
		Username:          "username",
		RdsInstanceRegion: "us-east-1",
		AwsCredentials: &common.AwsCredentials{
			AwsAccessKeyId:     pointer.To("AKIA000FAKE00KEY00ID"),
			AwsSecretAccessKey: pointer.To("Secret+Access+Key/0000000000000000000000"),
		},
	}

	tests := []struct {
		name  string
		given AwsRdsDatabaseServiceConfiguration
		want  error
	}{
		{
			name:  "database protocol is missing",
			given: AwsRdsDatabaseServiceConfiguration{
				// database protocol is missing
			},
			want: errors.New("database protocol is required"),
		},
		{
			name: "hostname is missing",
			given: AwsRdsDatabaseServiceConfiguration{
				DatabaseProtocol: testDatbaseProtocol,
				HostnameAndPort: HostnameAndPort{
					// hostname is missing
					Port: 3306,
				},
				AuthenticationType:      DatabaseAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuth: testPasswordAuth,
			},
			want: errors.New("hostname is a required field"),
		},
		{
			name: "invalid database authentication type",
			given: AwsRdsDatabaseServiceConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: "invalid",
			},
			want: errors.New("invalid database authentication type: invalid"),
		},
		{
			name: "when username and password auth is picked, iam auth config should be nil",
			given: AwsRdsDatabaseServiceConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				HostnameAndPort:         testHostnameAndPort,
				AuthenticationType:      DatabaseAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuth: testPasswordAuth,
				IamAuth:                 testIamAuth,
			},
			want: errors.New("authentication type is username_and_password, but IAM auth configuration is provided"),
		},
		{
			name: "username and password auth is picked, but username and password auth config is missing",
			given: AwsRdsDatabaseServiceConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: DatabaseAuthenticationTypeUsernameAndPassword,
				// username and password auth config is missing
			},
			want: errors.New("username and password auth configuration is required"),
		},
		{
			name: "when iam auth is picked, username and password auth config should be nil",
			given: AwsRdsDatabaseServiceConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				HostnameAndPort:         testHostnameAndPort,
				AuthenticationType:      DatabaseAuthenticationTypeIam,
				UsernameAndPasswordAuth: testPasswordAuth,
				IamAuth:                 testIamAuth,
			},
			want: errors.New("authentication type is iam, but username and password auth configuration is provided"),
		},
		{
			name: "iam auth is picked, but iam auth config is missing",
			given: AwsRdsDatabaseServiceConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: DatabaseAuthenticationTypeIam,
				// iam auth config is missing
			},
			want: errors.New("IAM auth configuration is required"),
		},
		{
			name: "happy path - username and password auth",
			given: AwsRdsDatabaseServiceConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				HostnameAndPort:         testHostnameAndPort,
				AuthenticationType:      DatabaseAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuth: testPasswordAuth,
			},
			want: nil,
		},
		{
			name: "happy path - iam auth",
			given: AwsRdsDatabaseServiceConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: DatabaseAuthenticationTypeIam,
				IamAuth:            testIamAuth,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_GcpCloudSqlDatabaseServiceConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testConnectorConfig := &GcpCloudSqlConnectorConfiguration{
		DatabaseProtocol:   DatabaseProtocolMySql,
		AuthenticationType: DatabaseAuthenticationTypeUsernameAndPassword,
		UsernameAndPasswordAuth: &GcpCloudSqlUsernameAndPasswordAuthConfiguration{
			UsernameAndPassword: UsernameAndPassword{
				Username: "username",
				Password: "password",
			},
			InstanceId:         "instance-id",
			GcpCredentialsJson: `{"something": "something"}`,
		},
	}
	testStandardConfig := &GcpCloudSqlStandardConfiguration{
		DatabaseProtocol: DatabaseProtocolMySql,
		HostnameAndPort: HostnameAndPort{
			Hostname: "hostname",
			Port:     3306,
		},
		AuthenticationType: DatabaseAuthenticationTypeTls,
		TlsAuth: &DatabaseTlsAuthConfiguration{
			UsernameAndPassword: UsernameAndPassword{
				Username: "username",
				Password: "password",
			},
			TlsConfig: TlsConfig{
				Certificate: "certificate",
				Key:         "private-key",
			},
		},
	}

	tests := []struct {
		name  string
		given GcpCloudSqlDatabaseServiceConfiguration
		want  error
	}{
		{
			name: "when connector is enabled, standard config should be nil",
			given: GcpCloudSqlDatabaseServiceConfiguration{
				CloudSqlConnectorEnabled: true,
				Connector:                testConnectorConfig,
				Standard:                 testStandardConfig,
			},
			want: errors.New("cloudsql_connector_enabled is true, but standard configuration is provided"),
		},
		{
			name: "connector is enabled, but connector config is missing",
			given: GcpCloudSqlDatabaseServiceConfiguration{
				CloudSqlConnectorEnabled: true,
				// connector config is missing
			},
			want: errors.New("Google Cloud SQL connector configuration is required"),
		},
		{
			name: "when connector is disabled, connector config should be nil",
			given: GcpCloudSqlDatabaseServiceConfiguration{
				CloudSqlConnectorEnabled: false,
				Connector:                testConnectorConfig,
				Standard:                 testStandardConfig,
			},
			want: errors.New("cloudsql_connector_enabled is false, but connector configuration is provided"),
		},
		{
			name: "connector is disabled, but standard config is missing",
			given: GcpCloudSqlDatabaseServiceConfiguration{
				CloudSqlConnectorEnabled: false,
				// standard config is missing
			},
			want: errors.New("standard Google Cloud SQL configuration is required"),
		},
		{
			name: "happy path - connector config",
			given: GcpCloudSqlDatabaseServiceConfiguration{
				CloudSqlConnectorEnabled: true,
				Connector:                testConnectorConfig,
			},
			want: nil,
		},
		{
			name: "happy path - standard config",
			given: GcpCloudSqlDatabaseServiceConfiguration{
				CloudSqlConnectorEnabled: false,
				Standard:                 testStandardConfig,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_GcpCloudSqlStandardConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testDatbaseProtocol := DatabaseProtocolMySql
	testHostnameAndPort := HostnameAndPort{
		Hostname: "hostname",
		Port:     3306,
	}
	testPasswordAuth := &DatabaseUsernameAndPasswordAuthConfiguration{
		UsernameAndPassword: UsernameAndPassword{
			Username: "username",
			Password: "password",
		},
	}
	testTlsAuth := &DatabaseTlsAuthConfiguration{
		UsernameAndPassword: UsernameAndPassword{
			Username: "username",
			Password: "password",
		},
		TlsConfig: TlsConfig{
			Certificate: "certificate",
			Key:         "private-key",
		},
	}

	tests := []struct {
		name  string
		given GcpCloudSqlStandardConfiguration
		want  error
	}{
		{
			name:  "database protocol is missing",
			given: GcpCloudSqlStandardConfiguration{
				// database protocol is missing
			},
			want: errors.New("database protocol is required"),
		},
		{
			name: "hostname is missing",
			given: GcpCloudSqlStandardConfiguration{
				DatabaseProtocol: testDatbaseProtocol,
				HostnameAndPort: HostnameAndPort{
					// hostname is missing
					Port: 3306,
				},
				AuthenticationType:      DatabaseAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuth: testPasswordAuth,
			},
			want: errors.New("hostname is a required field"),
		},
		{
			name: "invalid database authentication type",
			given: GcpCloudSqlStandardConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: "invalid",
			},
			want: errors.New("invalid database authentication type: invalid"),
		},
		{
			name: "when username and password auth is picked, tls auth config should be nil",
			given: GcpCloudSqlStandardConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				HostnameAndPort:         testHostnameAndPort,
				AuthenticationType:      DatabaseAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuth: testPasswordAuth,
				TlsAuth:                 testTlsAuth,
			},
			want: errors.New("authentication type is username_and_password, but TLS auth configuration is provided"),
		},
		{
			name: "username and password auth is picked, but username and password auth config is missing",
			given: GcpCloudSqlStandardConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: DatabaseAuthenticationTypeUsernameAndPassword,
				// username and password auth config is missing
			},
			want: errors.New("username and password auth configuration is required"),
		},
		{
			name: "when tls auth is picked, username and password auth config should be nil",
			given: GcpCloudSqlStandardConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				HostnameAndPort:         testHostnameAndPort,
				AuthenticationType:      DatabaseAuthenticationTypeTls,
				UsernameAndPasswordAuth: testPasswordAuth,
				TlsAuth:                 testTlsAuth,
			},
			want: errors.New("authentication type is tls, but username and password auth configuration is provided"),
		},
		{
			name: "tls auth is picked, but tls auth config is missing",
			given: GcpCloudSqlStandardConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: DatabaseAuthenticationTypeTls,
				// tls auth config is missing
			},
			want: errors.New("TLS auth configuration is required"),
		},
		{
			name: "happy path - username and password auth",
			given: GcpCloudSqlStandardConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				HostnameAndPort:         testHostnameAndPort,
				AuthenticationType:      DatabaseAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuth: testPasswordAuth,
			},
			want: nil,
		},
		{
			name: "happy path - tls auth",
			given: GcpCloudSqlStandardConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: DatabaseAuthenticationTypeTls,
				TlsAuth:            testTlsAuth,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_GcpCloudSqlConnectorConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testDatbaseProtocol := DatabaseProtocolMySql
	testPasswordAuth := &GcpCloudSqlUsernameAndPasswordAuthConfiguration{
		UsernameAndPassword: UsernameAndPassword{
			Username: "username",
			Password: "password",
		},
		InstanceId:         "instance-id",
		GcpCredentialsJson: `{"something": "something"}`,
	}
	testIamAuth := &GcpCloudSqlIamAuthConfiguration{
		Username:           "username",
		InstanceId:         "instance-id",
		GcpCredentialsJson: `{"something": "something"}`,
	}

	tests := []struct {
		name  string
		given GcpCloudSqlConnectorConfiguration
		want  error
	}{
		{
			name:  "database protocol is missing",
			given: GcpCloudSqlConnectorConfiguration{
				// database protocol is missing
			},
			want: errors.New("database protocol is required"),
		},
		{
			name: "invalid database authentication type",
			given: GcpCloudSqlConnectorConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				AuthenticationType: "invalid",
			},
			want: errors.New("invalid database authentication type: invalid"),
		},
		{
			name: "when username and password auth is picked, iam auth config should be nil",
			given: GcpCloudSqlConnectorConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				AuthenticationType:      DatabaseAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuth: testPasswordAuth,
				IamAuth:                 testIamAuth,
			},
			want: errors.New("authentication type is username_and_password, but IAM auth configuration is provided"),
		},
		{
			name: "username and password auth is picked, but username and password auth config is missing",
			given: GcpCloudSqlConnectorConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				AuthenticationType: DatabaseAuthenticationTypeUsernameAndPassword,
				// username and password auth config is missing
			},
			want: errors.New("username and password auth configuration is required"),
		},
		{
			name: "when iam auth is picked, username and password auth config should be nil",
			given: GcpCloudSqlConnectorConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				AuthenticationType:      DatabaseAuthenticationTypeIam,
				UsernameAndPasswordAuth: testPasswordAuth,
				IamAuth:                 testIamAuth,
			},
			want: errors.New("authentication type is iam, but username and password auth configuration is provided"),
		},
		{
			name: "iam auth is picked, but iam auth config is missing",
			given: GcpCloudSqlConnectorConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				AuthenticationType: DatabaseAuthenticationTypeIam,
				// iam auth config is missing
			},
			want: errors.New("IAM auth configuration is required"),
		},
		{
			name: "happy path - username and password auth",
			given: GcpCloudSqlConnectorConfiguration{
				DatabaseProtocol:        testDatbaseProtocol,
				AuthenticationType:      DatabaseAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuth: testPasswordAuth,
			},
			want: nil,
		},
		{
			name: "happy path - iam auth",
			given: GcpCloudSqlConnectorConfiguration{
				DatabaseProtocol:   testDatbaseProtocol,
				AuthenticationType: DatabaseAuthenticationTypeIam,
				IamAuth:            testIamAuth,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_DatabaseUsernameAndPasswordAuthConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testUsername := "username"
	testPassword := "password"

	tests := []struct {
		name  string
		given DatabaseUsernameAndPasswordAuthConfiguration
		want  error
	}{
		{
			name: "username is missing",
			given: DatabaseUsernameAndPasswordAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					// username is missing
					Password: testPassword,
				},
			},
			want: errors.New("username is required"),
		},
		{
			name: "password is missing",
			given: DatabaseUsernameAndPasswordAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					// password is missing
				},
			},
			want: errors.New("password is required"),
		},
		{
			name: "happy path",
			given: DatabaseUsernameAndPasswordAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					Password: testPassword,
				},
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_DatabaseTlsAuthConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testUsername := "username"
	testPassword := "password"
	testCertificate := "certificate"
	testKey := "private-key"

	tests := []struct {
		name  string
		given DatabaseTlsAuthConfiguration
		want  error
	}{
		{
			name: "username is missing",
			given: DatabaseTlsAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					// username is missing
					Password: testPassword,
				},
				TlsConfig: TlsConfig{
					Certificate: testCertificate,
					Key:         testKey,
				},
			},
			want: errors.New("username is required"),
		},
		{
			name: "password is missing",
			given: DatabaseTlsAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					// password is missing
				},
				TlsConfig: TlsConfig{
					Certificate: testCertificate,
					Key:         testKey,
				},
			},
			want: errors.New("password is required"),
		},
		{
			name: "tls certificate is missing",
			given: DatabaseTlsAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					Password: testPassword,
				},
				TlsConfig: TlsConfig{
					// certificate is missing
					Key: testKey,
				},
			},
			want: errors.New("TLS certificate is required"),
		},
		{
			name: "tls private key is missing",
			given: DatabaseTlsAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					Password: testPassword,
				},
				TlsConfig: TlsConfig{
					Certificate: testCertificate,
					// private key is missing
				},
			},
			want: errors.New("TLS private key is required"),
		},
		{
			name: "happy path",
			given: DatabaseTlsAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					Password: testPassword,
				},
				TlsConfig: TlsConfig{
					Certificate: testCertificate,
					Key:         testKey,
				},
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_AwsRdsUsernameAndPasswordAuthConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testUsername := "username"
	testPassword := "password"

	tests := []struct {
		name  string
		given AwsRdsUsernameAndPasswordAuthConfiguration
		want  error
	}{
		{
			name: "username is missing",
			given: AwsRdsUsernameAndPasswordAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					// username is missing
					Password: testPassword,
				},
			},
			want: errors.New("username is required"),
		},
		{
			name: "password is missing",
			given: AwsRdsUsernameAndPasswordAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					// password is missing
				},
			},
			want: errors.New("password is required"),
		},
		{
			name: "happy path",
			given: AwsRdsUsernameAndPasswordAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					Password: testPassword,
				},
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_AwsRdsIamAuthConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testRdsInstanceRegion := "us-east-1"
	testUsername := "username"
	testBadAwsCredentials := common.AwsCredentials{
		AwsAccessKeyId: pointer.To("AKIA000FAKE00KEY00ID"),
		// aws secret access key is missing
	}
	testGoodAwsCredentials := common.AwsCredentials{
		AwsAccessKeyId:     pointer.To("AKIA000FAKE00KEY00ID"),
		AwsSecretAccessKey: pointer.To("Secret+Access+Key/0000000000000000000000"),
	}

	tests := []struct {
		name  string
		given AwsRdsIamAuthConfiguration
		want  error
	}{
		{
			name: "aws rds instance region is missing",
			given: AwsRdsIamAuthConfiguration{
				Username: testUsername,
				// aws rds instance region is missing
				AwsCredentials: &testGoodAwsCredentials,
			},
			want: errors.New("AWS RDS instance region is required"),
		},
		{
			name: "username is missing",
			given: AwsRdsIamAuthConfiguration{
				// username is missing
				RdsInstanceRegion: testRdsInstanceRegion,
				AwsCredentials:    &testGoodAwsCredentials,
			},
			want: errors.New("username is required"),
		},
		{
			name: "bad aws credentials",
			given: AwsRdsIamAuthConfiguration{
				RdsInstanceRegion: testRdsInstanceRegion,
				Username:          testUsername,
				AwsCredentials:    &testBadAwsCredentials,
			},
			want: errors.New("invalid AWS credentials: aws_secret_access_key is required when aws_access_key_id is provided"),
		},
		{
			name: "happy path - with aws credentials",
			given: AwsRdsIamAuthConfiguration{
				RdsInstanceRegion: testRdsInstanceRegion,
				Username:          testUsername,
				AwsCredentials:    &testGoodAwsCredentials,
			},
			want: nil,
		},
		{
			name: "happy path - without aws credentials",
			given: AwsRdsIamAuthConfiguration{
				RdsInstanceRegion: testRdsInstanceRegion,
				Username:          testUsername,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			if test.want == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, test.want, err.Error())
			}
		})
	}
}

func Test_GcpCloudSqlUsernameAndPasswordAuthConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testUsername := "username"
	testPassword := "password"
	testInstanceId := "instance-id"
	testGcpCredentialsJson := `{"something": "something"}`

	tests := []struct {
		name  string
		given GcpCloudSqlUsernameAndPasswordAuthConfiguration
		want  error
	}{
		{
			name: "username is missing",
			given: GcpCloudSqlUsernameAndPasswordAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					// username is missing
					Password: testPassword,
				},
				InstanceId:         testInstanceId,
				GcpCredentialsJson: testGcpCredentialsJson,
			},
			want: errors.New("username is required"),
		},
		{
			name: "password is missing",
			given: GcpCloudSqlUsernameAndPasswordAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					// password is missing
				},
				InstanceId:         testInstanceId,
				GcpCredentialsJson: testGcpCredentialsJson,
			},
			want: errors.New("password is required"),
		},
		{
			name: "instance id is missing",
			given: GcpCloudSqlUsernameAndPasswordAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					Password: testPassword,
				},
				// instance id is missing
				GcpCredentialsJson: testGcpCredentialsJson,
			},
			want: errors.New("instance ID is required"),
		},
		{
			name: "gcp credentials json is missing",
			given: GcpCloudSqlUsernameAndPasswordAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					Password: testPassword,
				},
				InstanceId: testInstanceId,
				// gcp credentials json is missing
			},
			want: errors.New("GCP credentials JSON is required"),
		},
		{
			name: "happy path",
			given: GcpCloudSqlUsernameAndPasswordAuthConfiguration{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					Password: testPassword,
				},
				InstanceId:         testInstanceId,
				GcpCredentialsJson: testGcpCredentialsJson,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_GcpCloudSqlIamAuthConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testUsername := "username"
	testInstanceId := "instance-id"
	testGcpCredentialsJson := `{"something": "something"}`

	tests := []struct {
		name  string
		given GcpCloudSqlIamAuthConfiguration
		want  error
	}{
		{
			name: "username is missing",
			given: GcpCloudSqlIamAuthConfiguration{
				// username is missing
				InstanceId:         testInstanceId,
				GcpCredentialsJson: testGcpCredentialsJson,
			},
			want: errors.New("username is required"),
		},
		{
			name: "instance id is missing",
			given: GcpCloudSqlIamAuthConfiguration{
				Username: testUsername,
				// instance id is missing
				GcpCredentialsJson: testGcpCredentialsJson,
			},
			want: errors.New("instance ID is required"),
		},
		{
			name: "gcp credentials json is missing",
			given: GcpCloudSqlIamAuthConfiguration{
				Username:   testUsername,
				InstanceId: testInstanceId,
				// gcp credentials json is missing
			},
			want: errors.New("GCP credentials JSON is required"),
		},
		{
			name: "happy path",
			given: GcpCloudSqlIamAuthConfiguration{
				Username:           testUsername,
				InstanceId:         testInstanceId,
				GcpCredentialsJson: testGcpCredentialsJson,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}
