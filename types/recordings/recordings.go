package recordings

import "github.com/borderzero/border0-go/lib/types/set"

const (
	// RecordingTypeAsciinema is the recording type for recordings in https://asciinema.org/ format.
	RecordingTypeAsciinema = "asciinema"

	// RecordingTypeKubernetesAPIRequestLog is the recording type for kubernetes api events in table format.
	RecordingTypeKubernetesAPIRequestLog = "kubernetes_api_request_log"

	// RecordingTypeDatabaseQueryLog is the recording type for database queries in table format.
	RecordingTypeDatabaseQueryLog = "database_query_log"

	// RecordingTypeBrowserSessionStream is the recording type for http DOM snapshot recordings e.g. rrweb.
	RecordingTypeBrowserSessionStream = "browser_session_stream"

	// RecordingTypeHTTPAccessLog is the recording type for http access logs in table format.
	RecordingTypeHTTPAccessLog = "http_access_log"

	// RecordingTypeSnowflakeQueryLog is the recording type for snowflake queries in table format.
	RecordingTypeSnowflakeQueryLog = "snowflake_query_log"
)

// ValidRecordingTypes represents allowed values for recording types.
var ValidRecordingTypes = set.New(
	RecordingTypeAsciinema,
	RecordingTypeKubernetesAPIRequestLog,
	RecordingTypeDatabaseQueryLog,
	RecordingTypeBrowserSessionStream,
	RecordingTypeHTTPAccessLog,
	RecordingTypeSnowflakeQueryLog,
)
