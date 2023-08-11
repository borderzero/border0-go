package service

const (
	// DatabaseProtocolTypeMySql is the database
	// service protocol for mysql databases.
	DatabaseProtocolMySql = "mysql"

	// DatabaseServiceTypePostgreSql is the database
	// service protocol for postgresql databases.
	DatabaseProtocolPostgreSql = "postgresql"
)

const (
	// DatabaseAuthenticationTypeAwsIam is the authentication type
	// for databases that use AWS IAM credentials for authentication.
	DatabaseAuthenticationTypeAwsIam = "aws_iam"

	// DatabaseAuthenticationTypeUsernameAndPassword is the authentication type
	// for databases that use username and password for authentication.
	DatabaseAuthenticationTypeUsernameAndPassword = "username_and_password"
)

// DatabaseServiceConfiguration represents service
// configuration for database services (fka sockets).
type DatabaseServiceConfiguration struct {
	DatabaseServiceType string `json:"database_service_type"`

	// mutually exclusive fields below
	StandardDatabaseServiceConfiguration *StandardDatabaseServiceConfiguration `json:"standard_database_service_configuration,omitempty"`
	CloudSqlDatabaseServiceConfiguration *CloudSqlDatabaseServiceConfiguration `json:"cloudsql_database_service_configuration,omitempty"`
}

// StandardDatabaseServiceConfiguration represents service
// configuration for standard database services (fka sockets).
type StandardDatabaseServiceConfiguration struct {
	HostnameAndPort           // inherited
	DatabaseProtocol   string `json:"protocol"`
	AuthenticationType string `json:"authentication_type"`
}

// CloudSqlDatabaseServiceConfiguration
type CloudSqlDatabaseServiceConfiguration struct {
	UsingConnector bool `json:"using_cloudsql_connector"`

	CloudSqlConnectorConfiguration     *CloudSqlConnectorConfiguration     `json:"cloudsql_connector_configuration"`
	CloudSqlConnectorlessConfiguration *CloudSqlConnectorlessConfiguration `json:"cloudsql_connectorless_configuration"`
}

type CloudSqlConnectorConfiguration struct {
	CloudSqlInstanceId string `json:"cloudsql_instance_id"`
}

type CloudSqlConnectorlessConfiguration struct {
	HostnameAndPort // inherit
}
