package service

import (
	"errors"
	"fmt"

	"github.com/borderzero/border0-go/lib/types/nilcheck"
	"github.com/borderzero/border0-go/types/common"
)

// Database service types supported by Border0. Choose `standard` for self-managed databases.
// Use `aws_rds` for AWS RDS databases, and select `gcp_cloudsql` for Google Cloud SQL databases.
const (
	DatabaseServiceTypeStandard      = "standard"       // standard MySQL or PostgreSQL, supports TLS and password auth
	DatabaseServiceTypeAwsRds        = "aws_rds"        // AWS RDS database, supports IAM and password auth
	DatabaseServiceTypeGcpCloudSql   = "gcp_cloudsql"   // Google Cloud SQL database, supports IAM, TLS and password auth
	DatabaseServiceTypeAzureSql      = "azure_sql"      // Azure SQL database, supports SQL authentication, azure password auth
	DatabaseServiceTypeAwsDocumentDB = "aws_documentdb" // AWS DocumentDB database, supports TLS, password auth and IAM auth
)

const (
	// DatabaseProtocolTypeMySql is the database service protocol for mysql databases.
	DatabaseProtocolMySql = "mysql"

	// DatabaseServiceTypePostgres is the database service protocol for postgresql databases.
	DatabaseProtocolPostgres = "postgres"

	// DatabaseProtocolTypeMSSql is the database service protocol for mssql databases.
	DatabaseProtocolSqlserver = "mssql"

	// DatabaseProtocolTypeCockroachDB is the database service protocol for cockroachdb databases.
	DatabaseProtocolCockroachDB = "cockroachdb"

	// DatabaseProtocolTypeMongoDB is the database service protocol for mongodb databases.
	DatabaseProtocolMongoDB = "mongodb"
)

const (
	// DatabaseAuthenticationTypeIam is the authentication type for databases
	// that use IAM credentials for authentication.
	DatabaseAuthenticationTypeIam = "iam"

	// DatabaseAuthenticationTypeTls is the authentication type for databases
	// that use TLS certificates for authentication. When using this type of
	// authentication, the client must provide a TLS certificate and a TLS key.
	DatabaseAuthenticationTypeTls = "tls"

	// DatabaseAuthenticationTypeUsernameAndPassword is the authentication type
	// for databases that use username and password for authentication.
	DatabaseAuthenticationTypeUsernameAndPassword = "username_and_password"

	// DatabaseAuthenticationTypeSqlAuthentication is the authentication type
	// for databases that use SQL authentication for authentication.
	DatabaseAuthenticationTypeSqlAuthentication = "sql_authentication"

	// DatabaseAuthenticationTypeAzureADPassword is the authentication type
	// for databases that use Azure Active Directory with password for authentication.
	DatabaseAuthenticationTypeAzureADPassword = "azure_active_directory_password"

	// DatabaseAuthenticationTypeAzureADIntegrated is the authentication type
	// for databases that use Azure Active Directory Integrated for authentication.
	DatabaseAuthenticationTypeAzureADIntegrated = "azure_active_directory_integrated"

	// DatabaseAuthenticationTypeKerberos is the authentication type
	// for databases that use kerberos for authentication.
	DatabaseAuthenticationTypeKerberos = "kerberos"

	// DatabaseAuthenticationTypeNoAuth is the authentication type
	// for databases that do not require authentication.
	DatabaseAuthenticationTypeNoAuth = "no_auth"
)

// =======================================================================================
// Database service configuration schema
// - database service type: standard, aws_rds, gcp_cloudsql, azure_sql
// - standard (when database service type is standard)
//     - hostname and port
//     - database protocol: mysql, postgres
//     - authentication type: username_and_password, tls
//     - username and password auth (when authentication type is username_and_password)
//         - username
//         - password
//     - tls auth (when authentication type is tls)
//         - username
//         - password
//         - certificate
//         - key
//         - ca_certificate (optional)
// - sqlserver (when database service type is sqlserver)
//     - hostname and port
//     - authentication type: sql_server, kerberos
//     - sql_server auth (when authentication type is sql_server)
//         - username
//         - password
//         - ca_certificate (optional)
//     - kerberos auth (when authentication type is kerberos)
//         - username
//         - password
// - aws rds (when database service type is aws_rds)
//     - hostname and port
//     - database protocol: mysql, postgres
//     - authentication type: username_and_password, tls
//     - username and password auth (when authentication type is username_and_password)
//         - username
//         - password
//         - ca_certificate (optional)
//     - iam auth (when authentication type is iam)
//         - rds_instance_region
//         - aws credentials: aws_access_key_id, aws_secret_access_key, aws_session_token, aws_profile
//         - username
//         - ca_certificate (optional)
// - google cloud sql (when database service type is gcp_cloudsql)
//     - cloudsql_connector_enabled
//     - standard (when cloudsql_connector_enabled is false)
//         - hostname and port
//         - database protocol: mysql
//         - authentication type: username_and_password, tls
//         - username and password auth (when authentication type is username_and_password)
//             - username
//             - password
//         - tls auth (when authentication type is tls)
//             - username
//             - password
//             - ca_certificate
//             - certificate
//             - key
//     - connector (when cloudsql_connector_enabled is true)
//         - database protocol: mysql, postgres
//         - authentication type: username_password, iam
//	   - username and password auth (when authentication type is username_and_password)
//             - username
//             - password
//             - instance_id
//             - gcp_credentials_json
//	   - iam auth (when authentication type is iam)
// 	    	   - username
//             - instance_id
//             - gcp_credentials_json
// - azure sql (when database service type is azure_sql)
//     - hostname and port
//     - database protocol: mssql
//     - authentication type: sql_authentication, azure_active_directory_password, azure_active_directory_integrated, kerberos
//     - sql authentication (when authentication type is sql_authentication)
//         - username
//         - password
//     - azure active directory password (when authentication type is azure_active_directory_password)
//         - username
//         - password
//     - azure active directory integrated (when authentication type is azure_active_directory_integrated)
//     - kerberos (when authentication type is kerberos)
//         - username
//         - password
// =======================================================================================

// DatabaseServiceConfiguration represents service configuration for database services (aka sockets).
type DatabaseServiceConfiguration struct {
	DatabaseServiceType string `json:"database_service_type"`

	// mutually exclusive fields below
	Standard      *StandardDatabaseServiceConfiguration      `json:"standard_database_service_configuration,omitempty"`
	AwsRds        *AwsRdsDatabaseServiceConfiguration        `json:"aws_rds_database_service_configuration,omitempty"`
	GcpCloudSql   *GcpCloudSqlDatabaseServiceConfiguration   `json:"gcp_cloudsql_database_service_configuration,omitempty"`
	AzureSql      *AzureSqlDatabaseServiceConfiguration      `json:"azure_sql_database_service_configuration,omitempty"`
	AwsDocumentDB *AwsDocumentDBDatabaseServiceConfiguration `json:"aws_documentdb_database_service_configuration,omitempty"`
}

// Validate ensures that the `DatabaseServiceConfiguration` is valid.
func (config DatabaseServiceConfiguration) Validate() error {
	if config.DatabaseServiceType == "" {
		return errors.New("database service type is required")
	}
	switch config.DatabaseServiceType {
	case DatabaseServiceTypeStandard:
		if nilcheck.AnyNotNil(config.AwsRds, config.GcpCloudSql, config.AzureSql, config.AwsDocumentDB) {
			return errors.New("database service type standard can only have standard configuration defined")
		}
		if config.Standard == nil {
			return errors.New("standard database service configuration is required")
		}
		return config.Standard.Validate()
	case DatabaseServiceTypeAwsRds:
		if nilcheck.AnyNotNil(config.Standard, config.GcpCloudSql, config.AzureSql, config.AwsDocumentDB) {
			return errors.New("database service type aws_rds can only have aws rds configuration defined")
		}
		if config.AwsRds == nil {
			return errors.New("AWS RDS database service configuration is required")
		}
		return config.AwsRds.Validate()
	case DatabaseServiceTypeAwsDocumentDB:
		if nilcheck.AnyNotNil(config.Standard, config.GcpCloudSql, config.AzureSql, config.AwsRds) {
			return errors.New("database service type aws_documentdb can only have aws documentdb configuration defined")
		}
		if config.AwsDocumentDB == nil {
			return errors.New("AWS DocumentDB database service configuration is required")
		}
		return config.AwsDocumentDB.Validate()
	case DatabaseServiceTypeGcpCloudSql:
		if nilcheck.AnyNotNil(config.Standard, config.AwsRds, config.AzureSql, config.AwsDocumentDB) {
			return errors.New("database service type gcp_cloudsql can only have gcp cloudsql configuration defined")
		}
		if config.GcpCloudSql == nil {
			return errors.New("Google Cloud SQL database service configuration is required")
		}
		return config.GcpCloudSql.Validate()
	case DatabaseServiceTypeAzureSql:
		if nilcheck.AnyNotNil(config.Standard, config.AwsRds, config.GcpCloudSql, config.AwsDocumentDB) {
			return errors.New("database service type azure_sql can only have azure sql configuration defined")
		}
		if config.AzureSql == nil {
			return errors.New("Azure SQL database service configuration is required")
		}
		return config.AzureSql.Validate()
	}
	return fmt.Errorf("invalid database service type: %s", config.DatabaseServiceType)
}

// =======================================================================================
// Configurations for database services
// - standard
// - aws rds
// - google cloud sql
// - azure sql
// =======================================================================================

// StandardDatabaseServiceConfiguration represents service configuration for self-managed databases.
// Self-managed databases are databases that are not managed by a cloud provider. For example, a MySQL
// or PostgreSQL database running on your laptop, or in a VM running in your data center or in the cloud.
//
// Supported database protocols are: `mysql`, `postgres` and `mssql`. For upstream authentication, supported auth
// types are: `username_and_password` and `tls`.
type StandardDatabaseServiceConfiguration struct {
	HostnameAndPort

	DatabaseProtocol   string `json:"protocol"`
	AuthenticationType string `json:"authentication_type"`

	UsernameAndPasswordAuth *DatabaseUsernameAndPasswordAuthConfiguration `json:"username_and_password_auth_configuration,omitempty"`
	TlsAuth                 *DatabaseTlsAuthConfiguration                 `json:"tls_auth_configuration,omitempty"`
	Kerberos                *DatabaseKerberosAuthConfiguration            `json:"kerberos_configuration,omitempty"`
	SqlAuthentication       *DatabaseSqlAuthConfiguration                 `json:"sql_authentication_configuration,omitempty"`
}

// Validate ensures that the `StandardDatabaseServiceConfiguration` is valid.
func (config StandardDatabaseServiceConfiguration) Validate() error {
	if config.DatabaseProtocol == "" {
		return errors.New("database protocol is required")
	}

	if err := config.HostnameAndPort.Validate(); err != nil {
		return err
	}

	switch config.DatabaseProtocol {
	case DatabaseProtocolMySql, DatabaseProtocolPostgres, DatabaseProtocolCockroachDB:
		switch config.AuthenticationType {
		case DatabaseAuthenticationTypeUsernameAndPassword:
			if nilcheck.AnyNotNil(config.TlsAuth, config.Kerberos, config.SqlAuthentication) {
				return errors.New("authentication type is username_and_password, but tls_auth, kerberos or sql_authentication configuration is provided")
			}
			if config.UsernameAndPasswordAuth == nil {
				return errors.New("username and password auth configuration is required")
			}
			return config.UsernameAndPasswordAuth.Validate()
		case DatabaseAuthenticationTypeTls:
			if nilcheck.AnyNotNil(config.UsernameAndPasswordAuth, config.Kerberos, config.SqlAuthentication) {
				return errors.New("authentication type is tls, but username_and_password, kerberos or sql_authentication configuration is provided")
			}
			if config.TlsAuth == nil {
				return errors.New("TLS auth configuration is required")
			}
			return config.TlsAuth.Validate()
		default:
			return fmt.Errorf("invalid database authentication type: %s", config.AuthenticationType)
		}
	case DatabaseProtocolMongoDB:
		switch config.AuthenticationType {
		case DatabaseAuthenticationTypeUsernameAndPassword:
			if nilcheck.AnyNotNil(config.TlsAuth, config.Kerberos, config.SqlAuthentication) {
				return errors.New("authentication type is username_and_password, but tls_auth, kerberos or sql_authentication configuration is provided")
			}
			if config.UsernameAndPasswordAuth == nil {
				return errors.New("username and password auth configuration is required")
			}
			return config.UsernameAndPasswordAuth.Validate()
		case DatabaseAuthenticationTypeTls:
			if nilcheck.AnyNotNil(config.UsernameAndPasswordAuth, config.Kerberos, config.SqlAuthentication) {
				return errors.New("authentication type is tls, but username_and_password, kerberos or sql_authentication configuration is provided")
			}
			if config.TlsAuth == nil {
				return errors.New("TLS auth configuration is required")
			}
			return config.TlsAuth.Validate()
		case DatabaseAuthenticationTypeNoAuth:
			if nilcheck.AnyNotNil(config.UsernameAndPasswordAuth, config.TlsAuth, config.Kerberos, config.SqlAuthentication) {
				return errors.New("authentication type is no_auth, but username_and_password, tls_auth, kerberos or sql_authentication configuration is provided")
			}
			return nil // No auth configuration is required for no_auth authentication type
		default:
			return fmt.Errorf("invalid database authentication type: %s", config.AuthenticationType)
		}
	case DatabaseProtocolSqlserver:
		switch config.AuthenticationType {
		case DatabaseAuthenticationTypeKerberos:
			if nilcheck.AnyNotNil(config.TlsAuth, config.UsernameAndPasswordAuth, config.SqlAuthentication) {
				return errors.New("authentication type is kerberos, but username_and_password, tls_auth or sql_authentication configuration is provided")
			}
			if config.Kerberos == nil {
				return errors.New("kerberos configuration is required")
			}
			return config.Kerberos.Validate()
		case DatabaseAuthenticationTypeSqlAuthentication:
			if nilcheck.AnyNotNil(config.TlsAuth, config.UsernameAndPasswordAuth, config.Kerberos) {
				return errors.New("authentication type is sql_authentication, but username_and_password, tls_auth or kerberos configuration is provided")
			}
			if config.SqlAuthentication == nil {
				return errors.New("sql_authentication configuration is required")
			}
			return config.SqlAuthentication.Validate()
		default:
			return fmt.Errorf("invalid database authentication type: %s", config.AuthenticationType)
		}
	}

	return fmt.Errorf("invalid database protocol: %s", config.DatabaseProtocol)
}

// AwsRdsDatabaseServiceConfiguration represents service configuration for AWS RDS databases. AWS RDS databases
// are cloud managed MySQL or PostgreSQL databases.
//
// Supported database protocols are: `mysql` and `postgres`. For upstream authentication, supported auth types
// are: `username_password` and `iam`. When using IAM authentication, the client must provide AWS credentials,
// AWS region and a username. You can provide an optional CA certificate to verify the RDS database server's
// certificate.
type AwsRdsDatabaseServiceConfiguration struct {
	HostnameAndPort

	DatabaseProtocol   string `json:"protocol"`
	AuthenticationType string `json:"authentication_type"`

	UsernameAndPasswordAuth *AwsRdsUsernameAndPasswordAuthConfiguration `json:"username_and_password_auth_configuration,omitempty"`
	IamAuth                 *AwsRdsIamAuthConfiguration                 `json:"iam_auth_configuration,omitempty"`
}

// Validate ensures that the `AwsRdsDatabaseServiceConfiguration` is valid.
func (config AwsRdsDatabaseServiceConfiguration) Validate() error {
	if config.DatabaseProtocol == "" {
		return errors.New("database protocol is required")
	}

	if err := config.HostnameAndPort.Validate(); err != nil {
		return err
	}

	switch config.AuthenticationType {
	case DatabaseAuthenticationTypeUsernameAndPassword:
		if nilcheck.AnyNotNil(config.IamAuth) {
			return errors.New("authentication type is username_and_password, but IAM auth configuration is provided")
		}
		if config.UsernameAndPasswordAuth == nil {
			return errors.New("username and password auth configuration is required")
		}
		return config.UsernameAndPasswordAuth.Validate()
	case DatabaseAuthenticationTypeIam:
		if nilcheck.AnyNotNil(config.UsernameAndPasswordAuth) {
			return errors.New("authentication type is iam, but username and password auth configuration is provided")
		}
		if config.IamAuth == nil {
			return errors.New("IAM auth configuration is required")
		}
		return config.IamAuth.Validate()
	}
	return fmt.Errorf("invalid database authentication type: %s", config.AuthenticationType)
}

// AwsDocumentDBDatabaseServiceConfiguration represents service configuration for AWS DocumentDB databases. AWS DocumentDB databases
// are cloud managed MongoDB clusters.
//
// Supported database protocols are: `mongodb`. For upstream authentication, supported auth types
// are: `username_password` and `iam`. When using IAM authentication, the client must provide AWS credentials and
// AWS region . You can provide an optional CA certificate to verify the database server's
// certificate.
type AwsDocumentDBDatabaseServiceConfiguration struct {
	HostnameAndPort

	DatabaseProtocol   string `json:"protocol"`
	AuthenticationType string `json:"authentication_type"`

	UsernameAndPasswordAuth *AwsRdsUsernameAndPasswordAuthConfiguration `json:"username_and_password_auth_configuration,omitempty"`
	IamAuth                 *MongoDBAWSAuthConfiguration                `json:"iam_auth_configuration,omitempty"`
}

// Validate ensures that the `AwsDocumentDBDatabaseServiceConfiguration` is valid.
func (config AwsDocumentDBDatabaseServiceConfiguration) Validate() error {
	if config.DatabaseProtocol == "" {
		return errors.New("database protocol is required")
	}

	if config.DatabaseProtocol != DatabaseProtocolMongoDB {
		return fmt.Errorf("invalid database protocol for AWS DocumentDB: %s (only mongodb is supported)", config.DatabaseProtocol)
	}

	if err := config.HostnameAndPort.Validate(); err != nil {
		return err
	}

	switch config.AuthenticationType {
	case DatabaseAuthenticationTypeUsernameAndPassword:
		if config.IamAuth != nil {
			return errors.New("authentication type is username_and_password, but IAM auth configuration is provided")
		}
		if config.UsernameAndPasswordAuth == nil {
			return errors.New("username and password auth configuration is required")
		}
		return config.UsernameAndPasswordAuth.Validate()
	case DatabaseAuthenticationTypeIam:
		if config.UsernameAndPasswordAuth != nil {
			return errors.New("authentication type is iam, but username and password auth configuration is provided")
		}
		if config.IamAuth == nil {
			return errors.New("IAM auth configuration is required")
		}
		return config.IamAuth.Validate()
	}
	return fmt.Errorf("invalid database authentication type: %s", config.AuthenticationType)
}

// GcpCloudSqlDatabaseServiceConfiguration represents service configuration for Google Cloud SQL databases.
// Google Cloud SQL databases are cloud managed MySQL or PostgreSQL databases.
//
// Border0 currently supports two ways of connecting to Google Cloud SQL databases: with and without the Cloud SQL Connector.
// Use the corresponding configuration fields to configure the upstream connection.
type GcpCloudSqlDatabaseServiceConfiguration struct {
	HostnameAndPort
	DatabaseProtocol string `json:"protocol"`

	UsernameAndPasswordAuth     *DatabaseUsernameAndPasswordAuthConfiguration `json:"username_and_password_auth_configuration,omitempty"`
	TlsAuth                     *DatabaseTlsAuthConfiguration                 `json:"tls_auth_configuration,omitempty"`
	GcpCloudSQLConnectorAuth    *GcpCloudSqlConnectorAuthConfiguration        `json:"cloudsql_connector_configuration,omitempty"`
	GcpCloudSQLConnectorIAMAuth *GcpCloudSqlConnectorIamAuthConfiguration     `json:"cloudsql_connector_iam_configuration,omitempty"`
}

// Validate ensures that the `GcpCloudSqlDatabaseServiceConfiguration` is valid.
func (config GcpCloudSqlDatabaseServiceConfiguration) Validate() error {
	switch {
	case config.UsernameAndPasswordAuth != nil:
		if nilcheck.AnyNotNil(config.TlsAuth, config.GcpCloudSQLConnectorAuth, config.GcpCloudSQLConnectorIAMAuth) {
			return errors.New("authentication type is username_and_password_auth_configuration, but tls_auth_configuration, cloudsql_auth_configuration or cloudsql_iam_auth_configuration is provided")
		}

		if err := config.HostnameAndPort.Validate(); err != nil {
			return err
		}

		return config.UsernameAndPasswordAuth.Validate()
	case config.TlsAuth != nil:
		if nilcheck.AnyNotNil(config.UsernameAndPasswordAuth, config.GcpCloudSQLConnectorAuth, config.GcpCloudSQLConnectorIAMAuth) {
			return errors.New("authentication type is tls_auth_configuration, but username_and_password_auth_configuration, cloudsql_auth_configuration or cloudsql_iam_auth_configuration is provided")
		}

		if err := config.HostnameAndPort.Validate(); err != nil {
			return err
		}

		return config.TlsAuth.Validate()
	case config.GcpCloudSQLConnectorAuth != nil:
		if nilcheck.AnyNotNil(config.UsernameAndPasswordAuth, config.TlsAuth, config.GcpCloudSQLConnectorIAMAuth) {
			return errors.New("authentication type is cloudsql_auth_configuration, but username_and_password_auth_configuration, tls_auth_configuration or cloudsql_iam_auth_configuration is provided")
		}

		return config.GcpCloudSQLConnectorAuth.Validate()
	case config.GcpCloudSQLConnectorIAMAuth != nil:
		if nilcheck.AnyNotNil(config.UsernameAndPasswordAuth, config.TlsAuth, config.GcpCloudSQLConnectorAuth) {
			return errors.New("authentication type is cloudsql_iam_auth_configuration, but username_and_password_auth_configuration, tls_auth_configuration or cloudsql_auth_configuration is provided")
		}

		return config.GcpCloudSQLConnectorIAMAuth.Validate()
	default:
		return errors.New("one of the following authentication types is required: azure_active_directory_password, azure_active_directory_integrated, kerberos, sql_authentication")
	}
}

// GcpCloudSqlConnectorAuthConfiguration represents service configuration for Google Cloud SQL database that will
// be connected to the upstream using the Cloud SQL Connector.
type GcpCloudSqlConnectorAuthConfiguration struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	InstanceId         string `json:"instance_id"`
	GcpCredentialsJson string `json:"gcp_credentials_json"`
}

// Validate ensures that the `GcpCloudSqlConnectorAuthConfiguration` has all the required fields.
func (config GcpCloudSqlConnectorAuthConfiguration) Validate() error {
	if config.Username == "" {
		return errors.New("username is required")
	}
	if config.Password == "" {
		return errors.New("password is required")
	}
	if config.InstanceId == "" {
		return errors.New("instance ID is required")
	}

	return nil
}

// GcpCloudSqlConnectorIamAuthConfiguration represents service configuration for Google Cloud SQL database that will
// be connected to the upstream using the Cloud SQL Connector and IAM authentication.
type GcpCloudSqlConnectorIamAuthConfiguration struct {
	Username           string `json:"username"`
	InstanceId         string `json:"instance_id"`
	GcpCredentialsJson string `json:"gcp_credentials_json"`
}

// Validate ensures that the `GcpCloudSqlConnectorIamAuthConfiguration` has all the required fields.
func (config GcpCloudSqlConnectorIamAuthConfiguration) Validate() error {
	if config.Username == "" {
		return errors.New("username is required")
	}

	if config.InstanceId == "" {
		return errors.New("instance ID is required")
	}

	return nil
}

// AzureSqlDatabaseServiceConfiguration represents service configuration for Azure SQL Server databases.
//
// Border0 currently supports four ways of connecting to Azure SQL Server databases.
// Use the corresponding configuration fields to configure the upstream connection.
type AzureSqlDatabaseServiceConfiguration struct {
	HostnameAndPort
	DatabaseProtocol string `json:"protocol"`

	AzureActiveDirectoryPassword   *DatabaseUsernameAndPasswordAuthConfiguration `json:"azure_active_directory_password_configuration,omitempty"`
	AzureActiveDirectoryIntegrated *struct{}                                     `json:"azure_active_directory_integrated_configuration,omitempty"`
	Kerberos                       *DatabaseKerberosAuthConfiguration            `json:"kerberos_configuration,omitempty"`
	SqlAuthentication              *DatabaseSqlAuthConfiguration                 `json:"sql_authentication_configuration,omitempty"`
}

// Validate ensures that the `AzureSqlDatabaseServiceConfiguration` is valid.
func (config AzureSqlDatabaseServiceConfiguration) Validate() error {
	if err := config.HostnameAndPort.Validate(); err != nil {
		return err
	}

	switch {
	case config.AzureActiveDirectoryPassword != nil:
		if nilcheck.AnyNotNil(config.AzureActiveDirectoryIntegrated, config.Kerberos, config.SqlAuthentication) {
			return errors.New("authentication type is azure_active_directory_password_configuration, but azure_active_directory_integrated_configuration, kerberos_configuration or sql_authentication_configuration is provided")
		}

		return config.AzureActiveDirectoryPassword.Validate()
	case config.AzureActiveDirectoryIntegrated != nil:
		if nilcheck.AnyNotNil(config.AzureActiveDirectoryPassword, config.Kerberos, config.SqlAuthentication) {
			return errors.New("authentication type is azure_active_directory_integrated_configuration, but azure_active_directory_password_configuration, kerberos_configuration or sql_authentication_configuration is provided")
		}

		return nil
	case config.Kerberos != nil:
		if nilcheck.AnyNotNil(config.AzureActiveDirectoryPassword, config.AzureActiveDirectoryIntegrated, config.SqlAuthentication) {
			return errors.New("authentication type is kerberos_configuration, but azure_active_directory_password_configuration, azure_active_directory_integrated_configuration or sql_authentication_configuration is provided")
		}

		return config.Kerberos.Validate()
	case config.SqlAuthentication != nil:
		if nilcheck.AnyNotNil(config.AzureActiveDirectoryPassword, config.AzureActiveDirectoryIntegrated, config.Kerberos) {
			return errors.New("authentication type is sql_authentication_configuration, but azure_active_directory_password_configuration, azure_active_directory_integrated_configuration or kerberos_configuration is provided")
		}

		return config.SqlAuthentication.Validate()
	default:
		return errors.New("one of the following authentication types is required: azure_active_directory_password, azure_active_directory_integrated, kerberos, sql_authentication")
	}
}

// =======================================================================================
// Configurations for different database authentication types
// - standard:
//     - username_password: username, password
//     - tls: username, password, ca_certificate, certificate, key
// - aws_rds:
//     - username_password: username, password, ca_certificate (optional)
//     - iam: aws credentials, username, ca_certificate (optional)
// - gcp_cloudsql:
//     - without cloud sql connector:
//         - username_password: username, password
//         - tls: username, password, certificate, key
//     - with cloud sql connector:
//         - username_password: username, password, instance_id, credentials
//         - iam: username, instance_id, credentials
// =======================================================================================

// DatabaseUsernameAndPasswordAuthConfiguration represents auth configuration that based on username and password.
type DatabaseUsernameAndPasswordAuthConfiguration struct {
	UsernameAndPassword
}

// Validate ensures that the `DatabaseUsernameAndPasswordAuthConfiguration` has all the required fields.
func (config DatabaseUsernameAndPasswordAuthConfiguration) Validate() error {
	if config.Username == "" {
		return errors.New("username is required")
	}
	if config.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// DatabaseKerberosAuthConfiguration represents auth configuration that based on username and password.
type DatabaseKerberosAuthConfiguration struct {
	UsernameAndPassword
}

// Validate ensures that the `DatabaseKerberosAuthConfiguration` has all the required fields.
func (config DatabaseKerberosAuthConfiguration) Validate() error {
	if config.Username == "" {
		return errors.New("username is required")
	}
	if config.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// DatabaseSqlAuthConfiguration represents auth configuration that based on username and password.
type DatabaseSqlAuthConfiguration struct {
	UsernameAndPassword
}

// Validate ensures that the `DatabaseSqlAuthConfiguration` has all the required fields.
func (config DatabaseSqlAuthConfiguration) Validate() error {
	if config.Username == "" {
		return errors.New("username is required")
	}
	if config.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// DatabaseTlsAuthConfiguration represents auth configuration that uses TLS for securing the connection. You must
// provide a username, a password, a client certificate and a client key. Optionally you can provide a CA certificate
// to verify the server's certificate.
type DatabaseTlsAuthConfiguration struct {
	UsernameAndPassword
	TlsConfig
}

// Validate ensures that the `DatabaseTlsAuthConfiguration` has all the required fields.
func (config DatabaseTlsAuthConfiguration) Validate() error {
	if config.Username == "" {
		return errors.New("username is required")
	}
	if config.Password == "" {
		return errors.New("password is required")
	}
	if config.Certificate == "" && config.Key != "" {
		return errors.New("TLS certificate is required")
	}
	if config.Key == "" && config.Certificate != "" {
		return errors.New("TLS private key is required")
	}
	return nil
}

// AwsRdsUsernameAndPasswordAuthConfiguration represents auth configuration for AWS RDS databases that use username
// and password. Optionally you can provide AWS CA bundle to verify the server's certificate.
type AwsRdsUsernameAndPasswordAuthConfiguration struct {
	UsernameAndPassword
	CaCertificate string `json:"ca_certificate,omitempty"`
}

// Validate ensures that the `AwsRdsUsernameAndPasswordAuthConfiguration` has all the required fields.
func (config AwsRdsUsernameAndPasswordAuthConfiguration) Validate() error {
	if config.Username == "" {
		return errors.New("username is required")
	}
	if config.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// AwsRdsIamAuthConfiguration represents auth configuration for AWS RDS databases that use IAM authentication. You must
// provide AWS credentials and a username. Optionally AWS CA bundle can be supplied to verify the server's certificate.
type AwsRdsIamAuthConfiguration struct {
	AwsCredentials    *common.AwsCredentials `json:"aws_credentials,omitempty"`
	RdsInstanceRegion string                 `json:"rds_instance_region"`
	Username          string                 `json:"username"`
	CaCertificate     string                 `json:"ca_certificate,omitempty"`
}

// Validate ensures that the `AwsRdsIamAuthConfiguration` has the required field and that the AWS credentials are valid.
func (config AwsRdsIamAuthConfiguration) Validate() error {
	if config.RdsInstanceRegion == "" {
		return errors.New("AWS RDS instance region is required")
	}
	if config.Username == "" {
		return errors.New("username is required")
	}
	if config.AwsCredentials != nil {
		if err := config.AwsCredentials.Validate(); err != nil {
			return fmt.Errorf("invalid AWS credentials: %w", err)
		}
	}
	return nil
}

// MongoDBAWSAuthConfiguration represents auth configuration for AWS DocumentDB databases that use IAM authentication. You must
// provide AWS credentials. Optionally AWS CA bundle can be supplied to verify the server's certificate.
type MongoDBAWSAuthConfiguration struct {
	AwsCredentials *common.AwsCredentials `json:"aws_credentials,omitempty"`
	ClusterRegion  string                 `json:"cluster_region"`
	CaCertificate  string                 `json:"ca_certificate,omitempty"`
}

// Validate ensures that the `MongoDBAWSAuthConfiguration` has the required field and that the AWS credentials are valid.
func (config MongoDBAWSAuthConfiguration) Validate() error {
	if config.ClusterRegion == "" {
		return errors.New("AWS DocumentDB cluster region is required")
	}
	if config.AwsCredentials != nil {
		if err := config.AwsCredentials.Validate(); err != nil {
			return fmt.Errorf("invalid AWS credentials: %w", err)
		}
	}
	return nil
}
