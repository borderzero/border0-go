package service

// Database service types supported by Border0. Choose `standard` for self-managed databases.
// Use `aws_rds` for AWS RDS databases, and select `google_cloudsql` for Google Cloud SQL databases.
const (
	DatabaseServiceTypeStandard = "standard"        // standard MySQL or PostgreSQL, supports TLS and password auth
	DatabaseServiceTypeRds      = "aws_rds"         // AWS RDS database, supports IAM and password auth
	DatabaseServiceTypeCloudSql = "google_cloudsql" // Google Cloud SQL database, supports IAM, TLS and password auth
)

const (
	// DatabaseProtocolTypeMySql is the database service protocol for mysql databases.
	DatabaseProtocolMySql = "mysql"

	// DatabaseServiceTypePostgres is the database service protocol for postgresql databases.
	DatabaseProtocolPostgres = "postgres"
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
)

// =======================================================================================
// Database service configuration schema
// - database service type: standard, aws_rds, google_cloudsql
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
// - aws rds (when database service type is aws_rds)
//     - hostname and port
//     - database protocol: mysql, postgres
//     - authentication type: username_and_password, tls
//     - username and password auth (when authentication type is username_and_password)
//         - username
//         - password
//         - ca_certificate (optional)
//     - iam auth (when authentication type is iam)
//         - aws credentials: aws_access_key_id, aws_secret_access_key, aws_region, aws_session_token, aws_profile
//         - username
//         - ca_certificate (optional)
// - google cloud sql (when database service type is google_cloudsql)
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
//         - iam auth (when authentication type is iam)
// 	       - username
//             - instance_id
//             - gcp_credentials_json
// =======================================================================================

// DatabaseServiceConfiguration represents service configuration for database services (aka sockets).
type DatabaseServiceConfiguration struct {
	DatabaseServiceType string `json:"database_service_type"`

	// mutually exclusive fields below
	Standard       *StandardDatabaseServiceConfiguration       `json:"standard_database_service_configuration,omitempty"`
	AwsRds         *AwsRdsDatabaseServiceConfiguration         `json:"aws_rds_database_service_configuration,omitempty"`
	GoogleCloudSql *GoogleCloudSqlDatabaseServiceConfiguration `json:"google_cloudsql_database_service_configuration,omitempty"`
}

// =======================================================================================
// Configurations for database services
// - standard
// - aws rds
// - google cloud sql
// =======================================================================================

// StandardDatabaseServiceConfiguration represents service configuration for self-managed databases.
// Self-managed databases are databases that are not managed by a cloud provider. For example, a MySQL
// or PostgreSQL database running on your laptop, or in a VM running in your data center or in the cloud.
//
// Supported database protocols are: `mysql` and `postgres`. For upstream authentication, supported auth
// types are: `username_and_password` and `tls`.
type StandardDatabaseServiceConfiguration struct {
	HostnameAndPort

	DatabaseProtocol   string `json:"protocol"`
	AuthenticationType string `json:"authentication_type"`

	UsernameAndPasswordAuth *DatabaseUsernameAndPasswordAuthConfiguration `json:"username_and_password_auth_configuration,omitempty"`
	TlsAuth                 *DatabaseTlsAuthConfiguration                 `json:"tls_auth_configuration,omitempty"`
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

// GoogleCloudSqlDatabaseServiceConfiguration represents service configuration for Google Cloud SQL databases.
// Google Cloud SQL databases are cloud managed MySQL or PostgreSQL databases.
//
// Border0 currently supports two ways of connecting to Google Cloud SQL databases: with and without the Cloud SQL Connector.
// Use the corresponding configuration fields to configure the upstream connection.
type GoogleCloudSqlDatabaseServiceConfiguration struct {
	CloudSqlConnectorEnabled bool `json:"cloudsql_connector_enabled"`

	Standard  *GoogleCloudSqlStandardConfiguration  `json:"standard_configuration,omitempty"`
	Connector *GoogleCloudSqlConnectorConfiguration `json:"connector_configuration,omitempty"`
}

// =======================================================================================
// Configurations specifically made for Google Cloud SQL
// - standard: without cloud sql connector
// - connector: with cloud sql connector
// =======================================================================================

// GoogleCloudSqlStandardConfiguration represents service configuration for Google Cloud SQL databases that will
// be connected to the upstream _WITHOUT_ using the Cloud SQL Connector.
//
// Supported database protocol is: `mysql`. For upstream authentication, supported auth types are: `username_password`,
// and `tls`. When using TLS authentication, the client must provide a username, a password, a client certificate and a
// client key.
type GoogleCloudSqlStandardConfiguration struct {
	HostnameAndPort

	DatabaseProtocol   string `json:"protocol"`
	AuthenticationType string `json:"authentication_type"`

	UsernameAndPasswordAuth *DatabaseUsernameAndPasswordAuthConfiguration `json:"username_and_password_auth_configuration,omitempty"`
	TlsAuth                 *DatabaseTlsAuthConfiguration                 `json:"tls_auth_configuration,omitempty"`
}

// GoogleCloudSqlConnectorConfiguration represents service configuration for Google Cloud SQL databases that will be
// connected to the upstream using the Cloud SQL Connector.
//
// Supported database protocol is: `mysql`. For upstream authentication, supported auth types are: `username_password`,
// and `iam`. When using IAM authentication, the client must provide a username and an instance ID. You will need to
// supply google credentials that are copied from the JSON credentials file.
type GoogleCloudSqlConnectorConfiguration struct {
	DatabaseProtocol   string `json:"protocol"`
	AuthenticationType string `json:"authentication_type"`

	UsernameAndPasswordAuth *GoogleCloudSqlUsernameAndPasswordAuthConfiguration `json:"username_and_password_auth_configuration,omitempty"`
	IamAuth                 *GoogleCloudSqlIamAuthConfiguration                 `json:"iam_auth_configuration,omitempty"`
}

// =======================================================================================
// Configurations for different database authentication types
// - standard:
//     - username_password: username, password
//     - tls: username, password, ca_certificate, certificate, key
// - aws_rds:
//     - username_password: username, password, ca_certificate (optional)
//     - iam: aws credentials, username, ca_certificate (optional)
// - google_cloudsql:
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

// DatabaseTlsAuthConfiguration represents auth configuration that uses TLS for securing the connection. You must
// provide a username, a password, a client certificate and a client key. Optionally you can provide a CA certificate
// to verify the server's certificate.
type DatabaseTlsAuthConfiguration struct {
	UsernameAndPassword
	TlsConfig
}

// AwsRdsUsernameAndPasswordAuthConfiguration represents auth configuration for AWS RDS databases that use username
// and password. Optionally you can provide AWS CA bundle to verify the server's certificate.
type AwsRdsUsernameAndPasswordAuthConfiguration struct {
	UsernameAndPassword
	CaCertificate string `json:"ca_certificate,omitempty"`
}

// AwsRdsIamAuthConfiguration represents auth configuration for AWS RDS databases that use IAM authentication. You must
// provide AWS credentials and a username. Optionally AWS CA bundle can be supplied to verify the server's certificate.
type AwsRdsIamAuthConfiguration struct {
	AwsCredentials
	Username      string `json:"username"`
	CaCertificate string `json:"ca_certificate,omitempty"`
}

// GoogleCloudSqlUsernameAndPasswordAuthConfiguration represents auth configuration for Google Cloud SQL databases that
// use username and password for authentication, and are connected to the upstream using the Cloud SQL Connector.
// You must provide a username, a password, an Cloud SQL instance ID and Google credentials that are copied from the JSON
// credentials file.
type GoogleCloudSqlUsernameAndPasswordAuthConfiguration struct {
	UsernameAndPassword
	InstanceId         string `json:"instance_id"`
	GcpCredentialsJson string `json:"gcp_credentials_json"`
}

// GoogleCloudSqlIamAuthConfiguration represents auth configuration for Google Cloud SQL databases that use IAM authentication,
// and are connected to the upstream using the Cloud SQL Connector. You must provide a username, an Cloud SQL instance ID
// and Google credentials that are copied from the JSON credentials file.
type GoogleCloudSqlIamAuthConfiguration struct {
	Username           string `json:"username"`
	InstanceId         string `json:"instance_id"`
	GcpCredentialsJson string `json:"gcp_credentials_json"`
}
